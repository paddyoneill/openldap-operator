package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
	"github.com/nscaledev/openldap-operator/internal/builder"
)

var _ = Describe("Builder", func() {
	var Builder *builder.Builder
	var scheme *runtime.Scheme

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	})

	Context("create new builder", func() {
		It("should return a newbuilder with scheme", func() {
			Builder = builder.NewBuilder(scheme)
			Expect(Builder).ToNot(BeNil())
			Expect(Builder.Scheme).ToNot(BeNil())
			Expect(Builder.Scheme).To(Equal(scheme))
		})
	})
})
