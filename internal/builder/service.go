package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/paddyoneill/openldap-operator/api/v1alpha1"
)

func (builder *Builder) DirectoryService(directory *v1alpha1.Directory) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directory.ServiceName(),
			Namespace: directory.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "openldap",
				"app.kubernetes.io/instance":  directory.Name,
				"app.kubernetes.io/component": "directory",
			},
		},
	}

	service.Spec.Selector = map[string]string{
		"app.kubernetes.io/instance": directory.Name,
	}
	service.Spec.Type = directory.Spec.Service.Type
	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:        "ldap",
			Protocol:    corev1.ProtocolTCP,
			Port:        389,
			TargetPort:  intstr.FromInt(3389),
			AppProtocol: ptr.To("ldap"),
		},
	}

	return service, controllerutil.SetControllerReference(directory, service, builder.Scheme)
}
