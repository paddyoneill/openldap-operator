package builder

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/paddyoneill/openldap-operator/api/v1alpha1"
	"github.com/paddyoneill/openldap-operator/internal/utils"
)

func (builder *Builder) DirectorySecret(directory *v1alpha1.Directory) (*corev1.Secret, error) {
	password, err := utils.GenerateRandonPassword(24)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directory.SecretName(),
			Namespace: directory.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "openldap",
				"app.kubernetes.io/instance":  directory.Name,
				"app.kubernetes.io/component": "directory",
			},
		},
		Data: map[string][]byte{
			"password": password,
		},
	}

	return secret, controllerutil.SetControllerReference(directory, secret, builder.Scheme)
}
