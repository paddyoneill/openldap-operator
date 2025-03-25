package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1alpha1 "github.com/paddyoneill/openldap-operator/api/v1alpha1"
	"github.com/paddyoneill/openldap-operator/internal/builder"
)

var _ = Describe("Secret", func() {
	var scheme *runtime.Scheme
	var Builder *builder.Builder
	var directory *v1alpha1.Directory
	var secret *corev1.Secret
	var err error
	var expectedLabels map[string]string

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

		expectedLabels = map[string]string{
			"app.kubernetes.io/name":      "openldap",
			"app.kubernetes.io/instance":  "foo-directory",
			"app.kubernetes.io/component": "directory",
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

		It("adds expected labels", func() {
			Expect(secret.Labels).To(Equal(expectedLabels))
		})

		It("sets controller reference", func() {
			Expect(secret.ObjectMeta.OwnerReferences[0].Name).To(Equal("foo-directory"))
		})

		It("contains the correct number of keys", func() {
			Expect(len(secret.Data)).To(HaveLen(1))
		})

		It("contains a password", func() {
			Expect(secret.Data).Should(HaveKey("password"))
			Expect(secret.Data["password"]).ToNot(Equal([]byte{}))
		})
	})
})
