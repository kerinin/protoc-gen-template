package data

import (
	"fmt"
	"strings"

	"github.com/kerinin/protoc-gen-template/meta"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type fieldID string

// FieldSlice is a slice of fields
type FieldSlice []Field

type ScalarPath []Field
type ScalarPathSlice []ScalarPath

func (s ScalarPath) Scalar() Field { return s[len(s)-1] }

// Visible returns the visible values in the slice
func (s FieldSlice) Visible() FieldSlice {
	outputs := make([]Field, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s FieldSlice) NotDeprecated() FieldSlice {
	outputs := make([]Field, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Field describes a protobuf message field
type Field struct {
	idx         int
	id          fieldID
	data        *Data
	parent      messageID
	oneof       oneofID // Non-empty for oneof fields
	typeMessage messageID
	typeEnum    enumID

	Name     string
	Meta     meta.FieldMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options  descriptor.FieldOptions // Globally-defined field metadata
	Comments Comments
	Number   int32
	Label    descriptor.FieldDescriptorProto_Label
	Type     descriptor.FieldDescriptorProto_Type
	// For numeric types, contains the original text representation of the value.
	// For booleans, "true" or "false".
	// For strings, contains the default text contents (not escaped in any way).
	// For bytes, contains the C escaped value.  All bytes >= 128 are escaped.
	// TODO(kenton):  Base-64 encode?
	DefaultValue string
	// JSON name of this field. The value is set by protocol compiler. If the
	// user has set a "json_name" option on this field, that option's value
	// will be used. Otherwise, it's deduced from the field's name by converting
	// it to camelCase.
	JSONName string
}

func (f Field) String() string {
	return string(f.id)
}

// IsVisible returns true if the field's parent and type is visible and its
// visibility metadata `PUBLIC`
func (f Field) IsVisible() bool {
	if !f.Parent().IsVisible() {
		return false
	}
	if t := f.TypeMessage(); t != nil && !t.IsVisible() {
		return false
	}
	if t := f.TypeEnum(); t != nil && !t.IsVisible() {
		return false
	}
	return f.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the field's parent or type are deprecated or if
// its deprecation option is true
func (f Field) IsDeprecated() bool {
	if f.Parent().IsDeprecated() {
		return true
	}
	if t := f.TypeMessage(); t != nil && t.IsDeprecated() {
		return true
	}
	if t := f.TypeEnum(); t != nil && t.IsDeprecated() {
		return true
	}
	return f.Options.Deprecated != nil && *f.Options.Deprecated == true
}

// Parent returns the fields's parent message
func (f Field) Parent() Message {
	return *f.data.messages[f.parent]
}

// Oneof returns the oneof this field is a member of, else nil
func (f Field) Oneof() *Oneof {
	return f.data.oneofs[f.oneof]
}

// IsOneof returns true if the field is part of a oneof
func (f Field) IsOneof() bool {
	return f.Oneof() != nil
}

// TypeMessage returns the message type if the field is message-typed, else nil
func (f Field) TypeMessage() *Message {
	return f.data.messages[f.typeMessage]
}

// TypeEnum returns the enum type if the field is enum-typed
func (f Field) TypeEnum() *Enum {
	return f.data.enums[f.typeEnum]
}

// IsRepeated is true if the field's label is 'REPEATED'
func (f Field) IsRepeated() bool {
	return f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

// IsTypeDouble is true if the field's type is 'double'
func (f Field) IsTypeDouble() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_DOUBLE
}

// IsTypeFloat is true if the field's type is 'float'
func (f Field) IsTypeFloat() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_FLOAT
}

// IsTypeInt64 is true if the field's type is 'int64'
func (f Field) IsTypeInt64() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_INT64
}

// IsTypeUint64 is true if the field's type is 'uint64'
func (f Field) IsTypeUint64() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_UINT64
}

// IsTypeInt32 is true if the field's type is 'int32'
func (f Field) IsTypeInt32() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_INT32
}

// IsTypeFixed64 is true if the field's type is 'fixed64'
func (f Field) IsTypeFixed64() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_FIXED64
}

// IsTypeFixed32 is true if the field's type is 'fixed32'
func (f Field) IsTypeFixed32() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_FIXED32
}

// IsTypeBool is true if the field's type is 'bool'
func (f Field) IsTypeBool() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_BOOL
}

// IsTypeString is true if the field's type is 'string'
func (f Field) IsTypeString() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_STRING
}

// IsTypeBytes is true if the field's type is 'bytes'
func (f Field) IsTypeBytes() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_BYTES
}

// IsTypeUint32 is true if the field's type is 'uint32'
func (f Field) IsTypeUint32() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_UINT32
}

// IsTypeSfixed32 is true if the field's type is 'sfixed32'
func (f Field) IsTypeSfixed32() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_SFIXED32
}

// IsTypeSfixed64 is true if the field's type is 'sfixed64'
func (f Field) IsTypeSfixed64() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_SFIXED64
}

// IsTypeSint32 is true if the field's type is 'sint32'
func (f Field) IsTypeSint32() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_SINT32
}

// IsTypeSint64 is true if the field's type is 'sint64'
func (f Field) IsTypeSint64() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_SINT64
}

// IsTypeGroup is true if the field's type is 'group'
func (f Field) IsTypeGroup() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_GROUP
}

// IsTypeEnum is true if the field's type is an enum
func (f Field) IsTypeEnum() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM
}

// IsTypeMessage is true if the field's type is a message
func (f Field) IsTypeMessage() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

// TypeNameString returns a prettified string descripton of a field's type
//
// Examples:
//   {Type: "TYPE_STRING", TypeName: nil, Label: OPTIONAL}			-> "string"
//   {Type: "TYPE_STRING", TypeName: nil, Label: REQUIRED}			-> "string"
//   {Type: "TYPE_STRING", TypeName: nil, Label: REPEATED}			-> "[]string"
//   {Type: "TYPE_BYTES", TypeName: nil, Label: REPEATED}			-> "bytes"
//   {Type: "TYPE_MESSAGE", TypeName: ".pkg.Msg, Label: OPTIONAL}	-> "pkg.Msg"
//   {Type: "TYPE_ENUM", TypeName: ".pkg.Enm, Label: OPTIONAL}		-> "pkg.Enm"
//
func (f Field) TypeNameString() string {
	switch f.Label {
	case descriptor.FieldDescriptorProto_LABEL_REPEATED:
		return fmt.Sprintf("[]" + f.typeName())
	case descriptor.FieldDescriptorProto_LABEL_OPTIONAL:
		// NOTE: Prefix with "*" for syntax v2?
		return f.typeName()
	default:
		return f.typeName()
	}
}

func (f Field) typeName() string {
	if f.typeMessage != messageID("") {
		return strings.TrimPrefix(string(f.typeMessage), ".")
	}
	if f.typeEnum != enumID("") {
		return strings.TrimPrefix(string(f.typeEnum), ".")
	}

	switch f.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "double"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "group"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return "message"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "bytes"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "enum"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "sfixed32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "sfixed64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "sint32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "sint64"
	default:
		return f.Type.String()
	}
}

func newFieldMetadata(in *descriptor.FieldOptions) (out meta.FieldMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	ext := proto.GetExtension(in, meta.E_FieldMeta)
	return *ext.(*meta.FieldMetadata)
}

func derefFieldOptions(o *descriptor.FieldOptions) (_ descriptor.FieldOptions) {
	if o == nil {
		return
	}
	return *o
}

type sortedFieldsByIndex []Field

func (s sortedFieldsByIndex) Len() int           { return len(s) }
func (s sortedFieldsByIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedFieldsByIndex) Less(i, j int) bool { return s[i].idx < s[j].idx }

type sortedScalarsByIndex []ScalarPath

func (s sortedScalarsByIndex) Len() int           { return len(s) }
func (s sortedScalarsByIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedScalarsByIndex) Less(i, j int) bool { return s[i][0].idx < s[j][0].idx }
