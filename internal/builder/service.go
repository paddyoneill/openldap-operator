package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
)

//func NewService(obj client.Object) *corev1.Service {
//	return &corev1.Service{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      .ServiceName(),
//			Namespace: .Namespace,
//		},
//	}
//}

func (builder *Builder) DirectoryService(directory *v1alpha1.Directory) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directory.ServiceName(),
			Namespace: directory.Namespace,
		},
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
