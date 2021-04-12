package data

import (
	"sort"

	"github.com/ReturnPath/protoc-gen-template/meta"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type messageID string

// MessageSlice is a slice of messages
type MessageSlice []Message

// Visible returns the visible values in the slice
func (s MessageSlice) Visible() MessageSlice {
	outputs := make([]Message, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s MessageSlice) NotDeprecated() MessageSlice {
	outputs := make([]Message, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// ToGenerate returns the values in the slice defined in files to be generated
func (s MessageSlice) ToGenerate() MessageSlice {
	outputs := make([]Message, 0, len(s))
	for _, f := range s {
		if f.File().Generate {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotNested returns the non-nested values in the slice
func (s MessageSlice) NotNested() MessageSlice {
	outputs := make([]Message, 0, len(s))
	for _, f := range s {
		if !f.IsNested() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// Message describes a protobuf message definition
type Message struct {
	idx      int
	id       messageID
	data     *Data
	file     fileID
	parent   messageID // Non-empty for embedded messages
	fields   []fieldID
	messages []messageID
	enums    []enumID
	oneofs   []oneofID

	Name          string
	Meta          meta.MessageMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options       descriptor.MessageOptions // Globally-defined message metadata
	Comments      Comments
	ReservedTags  []descriptor.DescriptorProto_ReservedRange
	ReservedNames []string // Reserved field names, which may not be used by fields in the same message.
}

func (m Message) String() string {
	return string(m.id)
}

// IsVisible returns true if the message's file and any enclosing message are
// visible, and its visibility metadata is `PUBLIC`
func (m Message) IsVisible() bool {
	if !m.File().IsVisible() {
		return false
	}
	if t := m.Parent(); t != nil && !t.IsVisible() {
		return false
	}
	return m.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the message's file or any enclosing messages are
// deprecated, or if its deprecation option is true
func (m Message) IsDeprecated() bool {
	if m.File().IsDeprecated() {
		return true
	}
	if t := m.Parent(); t != nil && t.IsDeprecated() {
		return true
	}
	return m.Options.Deprecated != nil && *m.Options.Deprecated == true
}

// File returns the containing file
func (m Message) File() File {
	return *m.data.files[m.file]
}

// Parent returns the messages's parent when embedded in a message, else nil
func (m Message) Parent() *Message {
	return m.data.messages[m.parent]
}

// IsNested returns true if the message is embedded in another message
func (m Message) IsNested() bool {
	return m.parent != messageID("")
}

// Root returns the outermost ancestor of the message
func (m Message) Root() Message {
	if m.parent == messageID("") {
		return m
	}
	return m.Parent().Root()
}

// Fields returns a slice of the message's fields
func (m Message) Fields() FieldSlice {
	vs := make([]Field, 0, len(m.fields))
	for _, v := range m.fields {
		if m.data.fields[v].Type == descriptor.FieldDescriptorProto_TYPE_GROUP {
			continue
		}
		vs = append(vs, *m.data.fields[v])
	}
	sort.Sort(sortedFieldsByIndex(vs))
	return vs
}

// Messages returns a slice of nested messages
func (m Message) Messages() MessageSlice {
	vs := make([]Message, 0, len(m.messages))
	for _, v := range m.messages {
		vs = append(vs, *m.data.messages[v])
	}
	sort.Sort(sortedMessagesByIndex(vs))
	return vs
}

// Enums returns a slice of the message's enums
func (m Message) Enums() EnumSlice {
	vs := make([]Enum, 0, len(m.enums))
	for _, v := range m.enums {
		vs = append(vs, *m.data.enums[v])
	}
	return vs
}

// Oneofs returns a slice of the message's oneofs
func (m Message) Oneofs() OneofSlice {
	vs := make([]Oneof, 0, len(m.oneofs))
	for _, v := range m.oneofs {
		vs = append(vs, *m.data.oneofs[v])
	}
	return vs
}

func newMessageMetadata(in *descriptor.MessageOptions) (out meta.MessageMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	ext := proto.GetExtension(in, meta.E_MessageMeta)
	return *ext.(*meta.MessageMetadata)
}

func derefMessageOptions(o *descriptor.MessageOptions) (_ descriptor.MessageOptions) {
	if o == nil {
		return
	}
	return *o
}

type sortedMessagesByIndex []Message

func (s sortedMessagesByIndex) Len() int           { return len(s) }
func (s sortedMessagesByIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedMessagesByIndex) Less(i, j int) bool { return s[i].idx < s[j].idx }
