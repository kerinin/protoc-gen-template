package data

import (
	"github.com/kerinin/protoc-gen-template/meta"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type methodID string

// MethodSlice is a slice of methods
type MethodSlice []Method

// Visible returns the visible values in the slice
func (s MethodSlice) Visible() MethodSlice {
	outputs := make([]Method, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s MethodSlice) NotDeprecated() MethodSlice {
	outputs := make([]Method, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Method describes a protobuf service method
type Method struct {
	id         methodID
	data       *Data
	parent     serviceID
	inputType  messageID
	outputType messageID

	Name            string
	Meta            meta.MethodMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options         descriptor.MethodOptions // Globally-defined service metadata
	Comments        Comments
	ClientStreaming bool
	ServerStreaming bool
}

func (m Method) String() string {
	return string(m.id)
}

// IsVisible returns true if the method's service and its input/output messages
// are visible, and its visibility metadata is `PUBLIC`
func (m Method) IsVisible() bool {
	if !m.Parent().IsVisible() {
		return false
	}
	if !m.InputType().IsVisible() {
		return false
	}
	if !m.OutputType().IsVisible() {
		return false
	}
	return m.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the method's service or any of its input/output
// messages are deprecated, or its deprecation option is true
func (m Method) IsDeprecated() bool {
	if m.Parent().IsDeprecated() {
		return true
	}
	if m.InputType().IsDeprecated() {
		return true
	}
	if m.OutputType().IsDeprecated() {
		return true
	}
	return m.Options.Deprecated != nil && *m.Options.Deprecated == true
}

// Parent returns the method's parent service
func (m Method) Parent() Service {
	return *m.data.services[m.parent]
}

// InputType returns the method's input type
func (m Method) InputType() Message {
	return *m.data.messages[m.inputType]
}

// OutputType returns the method's output type
func (m Method) OutputType() Message {
	return *m.data.messages[m.outputType]
}

func newMethodMetadata(in *descriptor.MethodOptions) (out meta.MethodMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	ext := proto.GetExtension(in, meta.E_MethodMeta)
	return *ext.(*meta.MethodMetadata)
}

func derefMethodOptions(o *descriptor.MethodOptions) (_ descriptor.MethodOptions) {
	if o == nil {
		return
	}
	return *o
}
