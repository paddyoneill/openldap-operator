package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1alpha1 "github.com/paddyoneill/openldap-operator/api/v1alpha1"
	"github.com/paddyoneill/openldap-operator/internal/builder"
)

var _ = Describe("Statefulset", func() {
	var scheme *runtime.Scheme
	var Builder *builder.Builder
	var directory *v1alpha1.Directory
	var sts *appsv1.StatefulSet
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
				Image: "test-image:some-tag",
			},
		}
		sts, err = Builder.DirectoryStatefulSet(directory)
	})

	Context("create directory statefulset", func() {
		It("doesn't return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a statefulset with correct metadata", func() {
			Expect(sts).ToNot(BeNil())
			Expect(sts.Name).To(Equal(directory.StatefulSetName()))
			Expect(sts.Namespace).To(Equal("bar"))
		})

		It("adds expected labels", func() {
			Expect(sts.Labels).To(Equal(map[string]string{
				"app.kubernetes.io/name":      "openldap",
				"app.kubernetes.io/instance":  "foo-directory",
				"app.kubernetes.io/component": "directory",
			}))
		})

		It("sets controller reference", func() {
			Expect(sts.ObjectMeta.OwnerReferences[0].Name).To(Equal("foo-directory"))
		})

		It("sets selector label", func() {
			Expect(sts.Spec.Selector.MatchLabels).To(Equal(map[string]string{
				"app.kubernetes.io/instance": "foo-directory",
			}))
		})

		It("sets correct pod template labels", func() {
			Expect(sts.Spec.Template.ObjectMeta.Labels).To(Equal(map[string]string{
				"app.kubernetes.io/name":      "openldap",
				"app.kubernetes.io/instance":  "foo-directory",
				"app.kubernetes.io/component": "directory",
			}))
		})

		It("sets the correct number of replicas", func() {
			Expect(*sts.Spec.Replicas).To(Equal(int32(1)))
		})

		It("sets the correct service name", func() {
			Expect(sts.Spec.ServiceName).To(Equal("foo-directory"))
		})

		It("sets correct volumes", func() {
			Expect(len(sts.Spec.Template.Spec.Volumes)).To(Equal(2))
			Expect(sts.Spec.Template.Spec.Volumes).To(Equal([]corev1.Volume{
				{
					Name: "slapd-ldif",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "foo-directory-slapd-config",
						},
					},
				},
				{
					Name: "slapd-data-dir",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			}))
		})

		It("creates correct pod spec", func() {
			Expect(len(sts.Spec.Template.Spec.Containers)).To(Equal(1))
			Expect(sts.Spec.Template.Spec.Containers[0].Name).To(Equal("slapd"))
			Expect(sts.Spec.Template.Spec.Containers[0].Image).To(Equal("test-image:some-tag"))
			Expect(sts.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))
			Expect(sts.Spec.Template.Spec.Containers[0].Ports).To(Equal([]corev1.ContainerPort{
				{
					Name:          "ldap",
					ContainerPort: 389,
					Protocol:      corev1.ProtocolTCP,
				},
			}))
			Expect(sts.Spec.Template.Spec.Containers[0].VolumeMounts).To(Equal([]corev1.VolumeMount{
				{
					Name:      "slapd-ldif",
					ReadOnly:  true,
					MountPath: "/etc/openldap/slapd.ldif",
					SubPath:   "slapd_ldif",
				},
				{
					Name:      "slapd-data-dir",
					MountPath: "/etc/openldap/slapd.d",
				},
			}))
		})
	})
})
