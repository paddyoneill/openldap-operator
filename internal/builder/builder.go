package builder

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type Builder struct {
	*runtime.Scheme
}

func NewBuilder(scheme *runtime.Scheme) *Builder {
	return &Builder{
		Scheme: scheme,
	}
}
