/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
	"github.com/nscaledev/openldap-operator/internal/builder"
)

// DirectoryReconciler reconciles a Directory object
type DirectoryReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Builder  *builder.Builder
}

const (
	directoryFinalizer = "openldap.nscale.dev/directoryFinalizer"
)

// +kubebuilder:rbac:groups=openldap.nscale.com,resources=directories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openldap.nscale.com,resources=directories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=openldap.nscale.com,resources=directories/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile directory resource
func (r *DirectoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	directory := &v1alpha1.Directory{}
	if err := r.Get(ctx, req.NamespacedName, directory); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("directory not found, ignoring since it must have been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to retrieve directory")
		return ctrl.Result{}, err
	}

	// Set directory condition to Unknown if no condition is already defined
	if directory.Status.Conditions == nil || len(directory.Status.Conditions) == 0 {
		meta.SetStatusCondition(&directory.Status.Conditions, metav1.Condition{
			Type:    v1alpha1.DirectoryAvailableCondition,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciler for new directory",
		})
		if err := r.Status().Update(ctx, directory); err != nil {
			logger.Error(err, "Failed to update directory status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if needed
	if directory.ObjectMeta.DeletionTimestamp.IsZero() && !controllerutil.ContainsFinalizer(directory, directoryFinalizer) {
		controllerutil.AddFinalizer(directory, directoryFinalizer)
		if err := r.Update(ctx, directory); err != nil {
			logger.Error(err, "failed to update directory with finalizer")
			return ctrl.Result{}, err
		}
	}

	// Directory marked for deletion
	if !directory.ObjectMeta.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(directory, directoryFinalizer) {
			logger.Info("performing finalizer actions for directory")
			meta.SetStatusCondition(&directory.Status.Conditions, metav1.Condition{
				Type:    v1alpha1.DirectoryDegradedCondition,
				Status:  metav1.ConditionTrue,
				Reason:  "Finalizing",
				Message: "Performing finalizer actions",
			})
			if err := r.Status().Update(ctx, directory); err != nil {
				logger.Error(err, "failed to update directory status")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(directory, directoryFinalizer)
			if err := r.Update(ctx, directory); err != nil {
				logger.Error(err, "failed to remove finalizer from directory")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.reconcileSecret(ctx, directory); err != nil {
		logger.Error(err, "failed to reconcile directory secret")
		meta.SetStatusCondition(&directory.Status.Conditions, metav1.Condition{
			Type:    v1alpha1.DirectoryAvailableCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "Reconciling",
			Message: fmt.Sprintf("failed to create secret for directory %s: %s", directory.Name, err.Error()),
		})

		if err := r.Status().Update(ctx, directory); err != nil {
			logger.Error(err, "Failed to update directory status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, directory); err != nil {
		logger.Error(err, "failed to reconcile directory service")
		meta.SetStatusCondition(&directory.Status.Conditions, metav1.Condition{
			Type:    v1alpha1.DirectoryAvailableCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "Reconciling",
			Message: fmt.Sprintf("failed to create service for directory %s: %s", directory.Name, err.Error()),
		})

		if err := r.Status().Update(ctx, directory); err != nil {
			logger.Error(err, "Failed to update directory status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, err
	}

	meta.SetStatusCondition(&directory.Status.Conditions, metav1.Condition{
		Type:    v1alpha1.DirectoryAvailableCondition,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: "Successfully reconciled resources",
	})
	if err := r.Status().Update(ctx, directory); err != nil {
		logger.Error(err, "Failed to update directory status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DirectoryReconciler) reconcileSecret(ctx context.Context, directory *v1alpha1.Directory) error {
	desired, err := r.Builder.DirectorySecret(directory)
	if err != nil {
		return err
	}

	existing := &corev1.Secret{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(desired), existing); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return r.Create(ctx, desired)
	}

	slapdLdif, err := r.Builder.GenerateSlapdLdif(directory, existing.Data["password_hash"])
	if err != nil {
		return err
	}

	patch := client.MergeFrom(directory.DeepCopy())
	existing.Labels = desired.Labels
	existing.Data["slapd_ldif"] = slapdLdif

	return r.Patch(ctx, existing, patch)
}

func (r *DirectoryReconciler) reconcileService(ctx context.Context, directory *v1alpha1.Directory) error {
	desired, err := r.Builder.DirectoryService(directory)
	if err != nil {
		return err
	}

	existing := &corev1.Service{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(desired), existing); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return r.Create(ctx, desired)
	}

	patch := client.MergeFrom(existing.DeepCopy())
	existing.Labels = desired.Labels
	existing.Spec.Ports = desired.Spec.Ports
	existing.Spec.Type = desired.Spec.Type
	existing.Spec.Selector = desired.Spec.Selector

	return r.Patch(ctx, existing, patch)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DirectoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Directory{}).
		Named("directory").
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
