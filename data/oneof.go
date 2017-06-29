package data

import (
	"sort"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type oneofID string

// OneofSlice is a slice of oneofs
type OneofSlice []Oneof

// Visible returns the visible values in the slice
func (s OneofSlice) Visible() OneofSlice {
	outputs := make([]Oneof, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s OneofSlice) NotDeprecated() OneofSlice {
	outputs := make([]Oneof, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Oneof describes a protobuf message definition
type Oneof struct {
	id     oneofID
	data   *Data
	parent messageID
	fields []fieldID

	Name     string
	Options  descriptor.OneofOptions // Globally-defined message metadata
	Comments Comments
}

func (o Oneof) String() string {
	return string(o.id)
}

// IsVisible returns true if the oneof's parent message is visible
func (o Oneof) IsVisible() bool {
	if !o.Parent().IsVisible() {
		return false
	}
	// TODO: Oneof visibility metadata
	return true
}

// IsDeprecated returns true if the oneof's parent message is deprecated
func (o Oneof) IsDeprecated() bool {
	if o.Parent().IsDeprecated() {
		return true
	}
	return false
}

// Parent returns the method's parent service
func (o Oneof) Parent() Message {
	return *o.data.messages[o.parent]
}

// Fields returns a slice of the oneof's fields
func (o Oneof) Fields() FieldSlice {
	vs := make([]Field, 0, len(o.fields))
	for _, v := range o.fields {
		if o.data.fields[v].Type == descriptor.FieldDescriptorProto_TYPE_GROUP {
			continue
		}
		vs = append(vs, *o.data.fields[v])
	}
	sort.Sort(sortedFieldsByIndex(vs))
	return vs
}

func derefOneofOptions(o *descriptor.OneofOptions) (_ descriptor.OneofOptions) {
	if o == nil {
		return
	}
	return *o
}
