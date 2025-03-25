package builder

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/paddyoneill/openldap-operator/api/v1alpha1"
)

func (builder *Builder) DirectoryStatefulSet(directory *v1alpha1.Directory) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directory.StatefulSetName(),
			Namespace: directory.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "openldap",
				"app.kubernetes.io/instance":  directory.Name,
				"app.kubernetes.io/component": "directory",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/instance": directory.Name,
				},
			},
			ServiceName: directory.ServiceName(),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":      "openldap",
						"app.kubernetes.io/instance":  directory.Name,
						"app.kubernetes.io/component": "directory",
					},
				},
			},
		},
	}

	var replicas int32 = 1
	sts.Spec.Replicas = ptr.To(replicas)

	sts.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "slapd-data-dir",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	sts.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:            "slapd",
			Image:           directory.Spec.Image,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env: []corev1.EnvVar{
				{
					Name: "OPENLDAP_CN_CONFIG_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: directory.SecretName(),
							},
							Key: "password",
						},
					},
				},
				{
					Name:  "OPENLDAP_SCHEMAS",
					Value: directory.Spec.SlapdConfig.Schemas.Join(),
				},
			},
			Ports: []corev1.ContainerPort{
				{
					Name:          "ldap",
					ContainerPort: 389,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "slapd-data-dir",
					MountPath: "/etc/openldap/slapd.d",
				},
			},
		},
	}

	return sts, controllerutil.SetControllerReference(directory, sts, builder.Scheme)
}
