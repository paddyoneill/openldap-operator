//go:build !ignore_autogenerated

/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigDatabaseConfig) DeepCopyInto(out *ConfigDatabaseConfig) {
	*out = *in
	if in.Access != nil {
		in, out := &in.Access, &out.Access
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigDatabaseConfig.
func (in *ConfigDatabaseConfig) DeepCopy() *ConfigDatabaseConfig {
	if in == nil {
		return nil
	}
	out := new(ConfigDatabaseConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Directory) DeepCopyInto(out *Directory) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Directory.
func (in *Directory) DeepCopy() *Directory {
	if in == nil {
		return nil
	}
	out := new(Directory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Directory) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryList) DeepCopyInto(out *DirectoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Directory, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryList.
func (in *DirectoryList) DeepCopy() *DirectoryList {
	if in == nil {
		return nil
	}
	out := new(DirectoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DirectoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryServiceSpec) DeepCopyInto(out *DirectoryServiceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryServiceSpec.
func (in *DirectoryServiceSpec) DeepCopy() *DirectoryServiceSpec {
	if in == nil {
		return nil
	}
	out := new(DirectoryServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectorySpec) DeepCopyInto(out *DirectorySpec) {
	*out = *in
	if in.SlapdConfig != nil {
		in, out := &in.SlapdConfig, &out.SlapdConfig
		*out = new(SlapdConfigSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(DirectoryServiceSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectorySpec.
func (in *DirectorySpec) DeepCopy() *DirectorySpec {
	if in == nil {
		return nil
	}
	out := new(DirectorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DirectoryStatus) DeepCopyInto(out *DirectoryStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DirectoryStatus.
func (in *DirectoryStatus) DeepCopy() *DirectoryStatus {
	if in == nil {
		return nil
	}
	out := new(DirectoryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendDatabaseConfig) DeepCopyInto(out *FrontendDatabaseConfig) {
	*out = *in
	if in.Access != nil {
		in, out := &in.Access, &out.Access
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendDatabaseConfig.
func (in *FrontendDatabaseConfig) DeepCopy() *FrontendDatabaseConfig {
	if in == nil {
		return nil
	}
	out := new(FrontendDatabaseConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in SchemaList) DeepCopyInto(out *SchemaList) {
	{
		in := &in
		*out = make(SchemaList, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchemaList.
func (in SchemaList) DeepCopy() SchemaList {
	if in == nil {
		return nil
	}
	out := new(SchemaList)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SlapdConfigSpec) DeepCopyInto(out *SlapdConfigSpec) {
	*out = *in
	if in.Schemas != nil {
		in, out := &in.Schemas, &out.Schemas
		*out = make(SchemaList, len(*in))
		copy(*out, *in)
	}
	if in.Overlays != nil {
		in, out := &in.Overlays, &out.Overlays
		*out = make([]Overlay, len(*in))
		copy(*out, *in)
	}
	if in.FrontendDatabase != nil {
		in, out := &in.FrontendDatabase, &out.FrontendDatabase
		*out = new(FrontendDatabaseConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ConfigDatabase != nil {
		in, out := &in.ConfigDatabase, &out.ConfigDatabase
		*out = new(ConfigDatabaseConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SlapdConfigSpec.
func (in *SlapdConfigSpec) DeepCopy() *SlapdConfigSpec {
	if in == nil {
		return nil
	}
	out := new(SlapdConfigSpec)
	in.DeepCopyInto(out)
	return out
}
