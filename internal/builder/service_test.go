package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
	"github.com/nscaledev/openldap-operator/internal/builder"
)

var _ = Describe("DirectoryService", func() {
	var scheme *runtime.Scheme
	var Builder *builder.Builder
	var directory *v1alpha1.Directory
	var service *corev1.Service
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
			Spec: v1alpha1.DirectorySpec{
				Service: &v1alpha1.DirectoryServiceSpec{
					Type: "some-type",
				},
			},
		}
		service, err = Builder.DirectoryService(directory)
	})

	Context("creates directory service", func() {
		It("doesn't return an error", func() {
			Expect(err).To(BeNil())
		})

		It("returns a service", func() {
			Expect(service).ToNot(BeNil())
		})

		It("contains correct metadata", func() {
			Expect(service.Name).To(Equal(directory.Name))
			Expect(service.Namespace).To(Equal(directory.Namespace))
		})

		It("contains the correct ports", func() {
			Expect(len(service.Spec.Ports)).To(Equal(1))
			Expect(service.Spec.Ports).To(Equal([]corev1.ServicePort{
				{
					Name:        "ldap",
					Protocol:    corev1.ProtocolTCP,
					Port:        389,
					TargetPort:  intstr.FromInt(3389),
					AppProtocol: ptr.To("ldap"),
				},
			}))
		})

		It("sets service type from directory spec", func() {
			Expect(service.Spec.Type).To(Equal(directory.Spec.Service.Type))
		})

		It("sets controller reference", func() {
			Expect(service.ObjectMeta.OwnerReferences[0].Name).To(Equal(directory.Name))
		})
	})
})
