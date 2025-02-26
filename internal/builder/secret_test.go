package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
	"github.com/nscaledev/openldap-operator/internal/builder"
)

var _ = Describe("Secret", func() {
	var scheme *runtime.Scheme
	var Builder *builder.Builder
	var directory *v1alpha1.Directory
	var secret *corev1.Secret
	var err error

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		Builder = builder.NewBuilder(scheme)
		directory = &v1alpha1.Directory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-directory",
				Namespace: "bar",
			},
		}
		secret, err = Builder.DirectorySecret(directory)
	})

	Context("create directory secret", func() {
		It("doesn't return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a secret with correct metadata", func() {
			Expect(secret).ToNot(BeNil())
			Expect(secret.Name).To(Equal(directory.SecretName()))
			Expect(secret.Namespace).To(Equal(directory.Namespace))
		})

		It("sets controller reference", func() {
			Expect(secret.ObjectMeta.OwnerReferences[0].Name).To(Equal(directory.Name))
		})

		It("contains the correct number of keys", func() {
			Expect(len(secret.Data)).To(Equal(2))
		})

		It("contains a password", func() {
			Expect(secret.Data).Should(HaveKey("password"))
			Expect(secret.Data["password"]).ToNot(Equal([]byte{}))
		})

		It("contains a password_hash", func() {
			Expect(secret.Data).Should(HaveKey("password_hash"))
			Expect(secret.Data["password_hash"]).ToNot(Equal([]byte{}))
		})
	})
})
