package data

import (
	"github.com/kerinin/protoc-gen-template/meta"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type serviceID string

// ServiceSlice is a slice of services
type ServiceSlice []Service

// Visible returns the visible values in the slice
func (s ServiceSlice) Visible() ServiceSlice {
	outputs := make([]Service, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s ServiceSlice) NotDeprecated() ServiceSlice {
	outputs := make([]Service, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// ToGenerate returns the values in the slice defined in files to be generated
func (s ServiceSlice) ToGenerate() ServiceSlice {
	outputs := make([]Service, 0, len(s))
	for _, f := range s {
		if f.File().Generate {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Service describes a protobuf service
type Service struct {
	id      serviceID
	data    *Data
	file    fileID
	methods []methodID

	Name     string
	Meta     meta.ServiceMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options  descriptor.ServiceOptions // Globally-defined service metadata
	Comments Comments
}

func (s Service) String() string {
	return string(s.id)
}

// IsVisible returns true if the service's file is visible and its visibility
// metadata is `PUBLIC`
func (s Service) IsVisible() bool {
	if !s.File().IsVisible() {
		return false
	}
	return s.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the service's file is deprecated or its
// deprecation option is true
func (s Service) IsDeprecated() bool {
	if s.File().IsDeprecated() {
		return true
	}
	return s.Options.Deprecated != nil && *s.Options.Deprecated == true
}

// File returns the containing file
func (s Service) File() File {
	return *s.data.files[s.file]
}

// Methods returns a slice of the service's methods
func (s Service) Methods() MethodSlice {
	vs := make([]Method, 0, len(s.methods))
	for _, v := range s.methods {
		vs = append(vs, *s.data.methods[v])
	}
	return vs
}

func newServiceMetadata(in *descriptor.ServiceOptions) (out meta.ServiceMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	ext := proto.GetExtension(in, meta.E_ServiceMeta)
	return *ext.(*meta.ServiceMetadata)
}

func derefServiceOptions(o *descriptor.ServiceOptions) (_ descriptor.ServiceOptions) {
	if o == nil {
		return
	}
	return *o
}
