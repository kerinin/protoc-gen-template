package data

import (
	"github.com/ReturnPath/protoc-gen-template/meta"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumID string

// EnumSlice is a slice of enums
type EnumSlice []Enum

// Visible returns the visible values in the slice
func (s EnumSlice) Visible() EnumSlice {
	outputs := make([]Enum, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s EnumSlice) NotDeprecated() EnumSlice {
	outputs := make([]Enum, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// ToGenerate returns the values in the slice defined in files to be generated
func (s EnumSlice) ToGenerate() EnumSlice {
	outputs := make([]Enum, 0, len(s))
	for _, f := range s {
		if f.File().Generate {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotNested returns the non-nested values in the slice
func (s EnumSlice) NotNested() EnumSlice {
	outputs := make([]Enum, 0, len(s))
	for _, f := range s {
		if !f.IsNested() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Enum describes a protobuf enum
type Enum struct {
	id     enumID
	data   *Data
	file   fileID
	parent messageID // Non-empty for embedded enums
	values []enumValueID

	Name     string
	Meta     meta.EnumMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options  descriptor.EnumOptions // Globally-defined enum metadata
	Comments Comments
}

func (e Enum) String() string {
	return string(e.id)
}

// IsVisible returns true if the enumeration's file and any enclosing messages
// are visibile and its visibility metadata is `PUBLIC`
func (e Enum) IsVisible() bool {
	if !e.File().IsVisible() {
		return false
	}
	if t := e.Parent(); t != nil && !t.IsVisible() {
		return false
	}
	return e.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the enumeration's file or any enclosing messages
// are deprecated, or if its deprecation option is true
func (e Enum) IsDeprecated() bool {
	if e.File().IsVisible() {
		return true
	}
	if t := e.Parent(); t != nil && t.IsVisible() {
		return true
	}
	return e.Options.Deprecated != nil && *e.Options.Deprecated == true
}

// IsNested returns true if the enum is embedded in a message
func (e Enum) IsNested() bool {
	return e.parent != messageID("")
}

// File returns the containing file
func (e Enum) File() File {
	return *e.data.files[e.file]
}

// Parent returns the enum's parent when embedded in a message, else nil
func (e Enum) Parent() *Message {
	return e.data.messages[e.parent]
}

// Values returns a slice of the enum's values
func (e Enum) Values() EnumValueSlice {
	vs := make([]EnumValue, 0, len(e.values))
	for _, v := range e.values {
		vs = append(vs, *e.data.enumValues[v])
	}
	return vs
}

func newEnumMetadata(in *descriptor.EnumOptions) (out meta.EnumMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	if ext, err := proto.GetExtension(in, meta.E_EnumMeta); err == nil {
		return *ext.(*meta.EnumMetadata)
	}

	return
}

func derefEnumOptions(o *descriptor.EnumOptions) (_ descriptor.EnumOptions) {
	if o == nil {
		return
	}
	return *o
}
