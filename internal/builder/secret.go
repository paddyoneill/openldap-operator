package builder

import (
	"bytes"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1alpha1 "github.com/nscaledev/openldap-operator/api/v1alpha1"
	"github.com/nscaledev/openldap-operator/internal/utils"
)

func (builder *Builder) DirectorySecret(directory *v1alpha1.Directory) (*corev1.Secret, error) {
	password, err := utils.GenerateRandonPassword(24)
	if err != nil {
		return nil, err
	}

	passwordHash, err := utils.Argon2HashPassword(password)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      directory.SecretName(),
			Namespace: directory.Namespace,
		},
		Data: map[string][]byte{
			"password":      password,
			"password_hash": passwordHash,
			"slapd_ldif":    []byte{},
		},
	}

	return secret, controllerutil.SetControllerReference(directory, secret, builder.Scheme)
}

func (builder *Builder) GenerateSlapdLdif(directory *v1alpha1.Directory, rootPwHash []byte) ([]byte, error) {
	var renderedConfig bytes.Buffer
	tmpl, err := template.New("slapd-ldif").Parse(slapd_ldif_tempalte)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(&renderedConfig, struct {
		*v1alpha1.SlapdConfigSpec
		RootPwHash string
	}{
		SlapdConfigSpec: directory.Spec.SlapdConfig,
		RootPwHash:      string(rootPwHash),
	}); err != nil {
		return nil, err
	}

	return renderedConfig.Bytes(), nil
}

var slapd_ldif_tempalte = `
dn: cn=config
objectClass: olcGlobal
cn: config

dn: cn=module,cn=config
objectClass: olcModuleList
cn: module
olcModulePath: /usr/lib/openldap
olcModuleLoad: argon2.so

dn: cn=schema,cn=config
objectClass: olcSchemaConfig
cn: schema
{{ range .Schemas }}
include: file:///etc/openldap/schema/{{ . }}.ldif
{{- end }}

dn: olcDatabase=frontend,cn=config
objectClass: olcDatabaseConfig
olcDatabase: frontend
{{- range .FrontendDatabase.Access }}
olcAccess: {{ . }}
{{- end }}

dn: olcDatabase=config,cn=config
objectClass: olcDatabaseConfig
olcDatabase: config
olcRootPW: {{ .RootPwHash }}
{{- range .ConfigDatabase.Access }}
olcAccess: {{ . }}
{{- end }}
`
