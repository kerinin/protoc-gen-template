package data

import (
	"github.com/ReturnPath/protoc-gen-template/meta"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type enumValueID string

// EnumValueSlice is a slice of enum values
type EnumValueSlice []EnumValue

// Visible returns the visible values in the slice
func (s EnumValueSlice) Visible() EnumValueSlice {
	outputs := make([]EnumValue, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s EnumValueSlice) NotDeprecated() EnumValueSlice {
	outputs := make([]EnumValue, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// EnumValue describes a protobuf enum value
type EnumValue struct {
	id     enumValueID
	data   *Data
	parent enumID

	Name     string
	Meta     meta.EnumValueMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options  descriptor.EnumValueOptions // Globally-defined enum value metadata
	Comments Comments
	Number   int32
}

func (e EnumValue) String() string {
	return string(e.id)
}

// IsVisible returns true if the enumeration value's parent is visible and its
// visibility metadata is `PUBLIC`
func (e EnumValue) IsVisible() bool {
	if !e.Parent().IsVisible() {
		return false
	}
	return e.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the enumeration value's parent is deprecated or its
// deprecation option is true
func (e EnumValue) IsDeprecated() bool {
	if e.Parent().IsDeprecated() {
		return true
	}
	return e.Options.Deprecated != nil && *e.Options.Deprecated == true
}

// Parent returns the enum for which this is a value
func (e EnumValue) Parent() Enum {
	return *e.data.enums[e.parent]
}

func newEnumValueMetadata(in *descriptor.EnumValueOptions) (out meta.EnumValueMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	if ext, err := proto.GetExtension(in, meta.E_EnumValueMeta); err == nil {
		return *ext.(*meta.EnumValueMetadata)
	}

	return
}

func derefEnumValueOptions(o *descriptor.EnumValueOptions) (_ descriptor.EnumValueOptions) {
	if o == nil {
		return
	}
	return *o
}
