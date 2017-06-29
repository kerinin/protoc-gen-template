package data

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ReturnPath/protoc-gen-template/meta"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/kr/pretty"
)

var request plugin.CodeGeneratorRequest

func TestMain(m *testing.M) {
	requestBytes, err := ioutil.ReadFile("testdata/dump.pb")
	if err != nil {
		log.Fatal("failed to read testdata/dump.pb: %s", err)
	}

	err = proto.Unmarshal(requestBytes, &request)
	if err != nil {
		log.Fatal("failed to unmarshal testdata/dump.pb: %s", err)
	}

	os.Exit(m.Run())
}

func stringPointer(s string) *string {
	return &s
}
func boolPointer(b bool) *bool {
	return &b
}

func TestFiles(t *testing.T) {
	var (
		expected = []File{
			File{
				id: ".testv2:protoc-gen-template/data/testdata/testv2.proto",
				messages: []messageID{
					messageID(".testv2.Message"),
					messageID(".testv2.OtherMessage"),
				},
				enums: []enumID{
					enumID(".testv2.Enum"),
					enumID(".testv2.OtherEnum"),
				},
				services: []serviceID{
					serviceID(".testv2.Service"),
					serviceID(".testv2.OtherService"),
				},
				Package: "testv2",
				Name:    "protoc-gen-template/data/testdata/testv2.proto",
				Comments: Comments{
					Leading:         " package comment\n",
					Trailing:        "",
					LeadingDetached: []string{},
				},
				Meta: meta.FileMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.FileOptions{
					GoPackage: stringPointer("testv2"),
				},
				Generate:     true,
				Dependencies: []string{"src/template/meta.proto"},
				Syntax:       "proto2",
			},
			File{
				id: ".testv3:protoc-gen-template/data/testdata/testv3.proto",
				messages: []messageID{
					messageID(".testv3.Message"),
					messageID(".testv3.OtherMessage"),
				},
				enums: []enumID{
					enumID(".testv3.Enum"),
					enumID(".testv3.OtherEnum"),
				},
				services: []serviceID{
					serviceID(".testv3.Service"),
					serviceID(".testv3.OtherService"),
				},
				Package: "testv3",
				Name:    "protoc-gen-template/data/testdata/testv3.proto",
				Comments: Comments{
					Leading:         " package comment\n",
					Trailing:        "",
					LeadingDetached: []string{},
				},
				Meta: meta.FileMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.FileOptions{
					GoPackage: stringPointer("testv3"),
				},
				Generate:     true,
				Dependencies: []string{"src/template/meta.proto"},
				Syntax:       "proto3",
			},
		}
		actual = New(&request).ToGenerate().Files()
	)

	for i := 0; i < len(actual); i++ {
		testDiff(t, "id", expected[i].id, actual[i].id)
		testDiff(t, "messages", expected[i].messages, actual[i].messages)
		testDiff(t, "enums", expected[i].enums, actual[i].enums)
		testDiff(t, "services", expected[i].services, actual[i].services)
		testDiff(t, "Package", expected[i].Package, actual[i].Package)
		testDiff(t, "Name", expected[i].Name, actual[i].Name)
		testDiff(t, "Comments", expected[i].Comments, actual[i].Comments)
		testDiff(t, "Meta", expected[i].Meta, actual[i].Meta)
		testDiff(t, "Options.GoPackage", expected[i].Options.GoPackage, actual[i].Options.GoPackage)
		testDiff(t, "Generate", expected[i].Generate, actual[i].Generate)
		testDiff(t, "Dependencies", expected[i].Dependencies, actual[i].Dependencies)
		testDiff(t, "Syntax", expected[i].Syntax, actual[i].Syntax)
	}
}

func TestMessages(t *testing.T) {
	var (
		expected = []Message{
			Message{
				id:     ".testv2.Message",
				file:   fileID(".testv2:protoc-gen-template/data/testdata/testv2.proto"),
				parent: messageID(""),
				fields: []fieldID{
					fieldID(".testv2.Message:string_field"),
					fieldID(".testv2.Message:repeated_string_field"),
					fieldID(".testv2.Message:enum_field"),
					fieldID(".testv2.Message:other_message_field"),
					fieldID(".testv2.Message:bool_field"),
					fieldID(".testv2.Message:embedded_enum_field"),
					fieldID(".testv2.Message:embedded_message_field"),
				},
				messages: []messageID{
					messageID(".testv2.Message.EmbeddedMessage"),
					messageID(".testv2.Message.OtherEmbeddedMessage"),
				},
				enums: []enumID{
					enumID(".testv2.Message.EmbeddedEnum"),
					enumID(".testv2.Message.OtherEmbeddedEnum"),
				},
				oneofs: []oneofID{
					oneofID(".testv2.Message:oneof_field"),
				},
				Name: "Message",
				Meta: meta.MessageMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.MessageOptions{
					Deprecated: boolPointer(true),
				},
				Comments: Comments{
					Leading:         " Message comment\n",
					Trailing:        "",
					LeadingDetached: []string{},
				},
				ReservedTags:  []descriptor.DescriptorProto_ReservedRange{},
				ReservedNames: []string{},
			},
			Message{
				id:     ".testv2.Message.EmbeddedMessage",
				file:   fileID(".testv2:protoc-gen-template/data/testdata/testv2.proto"),
				parent: messageID(".testv2.Message"),
				fields: []fieldID{
					fieldID(".testv2.Message.EmbeddedMessage:string_field"),
					fieldID(".testv2.Message.EmbeddedMessage:uint32_field"),
				},
				Name: "EmbeddedMessage",
				Comments: Comments{
					Leading: " EmbeddedMessage comment\n",
				},
			},
			Message{
				id:     ".testv2.Message.OtherEmbeddedMessage",
				file:   fileID(".testv2:protoc-gen-template/data/testdata/testv2.proto"),
				parent: messageID(".testv2.Message"),
				Name:   "OtherEmbeddedMessage",
			},
			Message{
				id:     ".testv2.OtherMessage",
				file:   fileID(".testv2:protoc-gen-template/data/testdata/testv2.proto"),
				parent: messageID(""),
				Name:   "OtherMessage",
			},
			Message{
				id:     ".testv3.Message",
				file:   fileID(".testv3:protoc-gen-template/data/testdata/testv3.proto"),
				parent: messageID(""),
				fields: []fieldID{
					fieldID(".testv3.Message:string_field"),
					fieldID(".testv3.Message:repeated_string_field"),
					fieldID(".testv3.Message:enum_field"),
					fieldID(".testv3.Message:other_message_field"),
					fieldID(".testv3.Message:bool_field"),
					fieldID(".testv3.Message:embedded_enum_field"),
					fieldID(".testv3.Message:embedded_message_field"),
				},
				messages: []messageID{
					messageID(".testv3.Message.EmbeddedMessage"),
					messageID(".testv3.Message.OtherEmbeddedMessage"),
				},
				enums: []enumID{
					enumID(".testv3.Message.EmbeddedEnum"),
					enumID(".testv3.Message.OtherEmbeddedEnum"),
				},
				oneofs: []oneofID{
					oneofID(".testv3.Message:oneof_field"),
				},
				Name: "Message",
				Meta: meta.MessageMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.MessageOptions{
					Deprecated: boolPointer(true),
				},
				Comments: Comments{
					Leading:         " Message comment\n",
					Trailing:        "",
					LeadingDetached: []string{},
				},
				ReservedTags:  []descriptor.DescriptorProto_ReservedRange{},
				ReservedNames: []string{},
			},
			Message{
				id:     ".testv3.Message.EmbeddedMessage",
				file:   fileID(".testv3:protoc-gen-template/data/testdata/testv3.proto"),
				parent: messageID(".testv3.Message"),
				fields: []fieldID{
					fieldID(".testv3.Message.EmbeddedMessage:string_field"),
					fieldID(".testv3.Message.EmbeddedMessage:uint32_field"),
				},
				Name: "EmbeddedMessage",
				Comments: Comments{
					Leading: " EmbeddedMessage comment\n",
				},
			},
			Message{
				id:     ".testv3.Message.OtherEmbeddedMessage",
				file:   fileID(".testv3:protoc-gen-template/data/testdata/testv3.proto"),
				parent: messageID(".testv3.Message"),
				Name:   "OtherEmbeddedMessage",
			},
			Message{
				id:     ".testv3.OtherMessage",
				file:   fileID(".testv3:protoc-gen-template/data/testdata/testv3.proto"),
				parent: messageID(""),
				Name:   "OtherMessage",
			},
		}
		actual = New(&request).ToGenerate().Messages()
	)

	for i := 0; i < len(actual); i++ {
		// for i := 0; i < 3; i++ {
		testDiff(t, "id", expected[i].id, actual[i].id)
		testDiff(t, "file", expected[i].file, actual[i].file)
		testDiff(t, "parent", expected[i].parent, actual[i].parent)
		testDiff(t, "fields", expected[i].fields, actual[i].fields)
		testDiff(t, "messages", expected[i].messages, actual[i].messages)
		testDiff(t, "enums", expected[i].enums, actual[i].enums)
		testDiff(t, "oneofs", expected[i].oneofs, actual[i].oneofs)
		testDiff(t, "Name", expected[i].Name, actual[i].Name)
		testDiff(t, "Meta", expected[i].Meta, actual[i].Meta)
		// testDiff(t, "Options.Deprecated", expected[i].Options.Deprecated, actual[i].Options.Deprecated)
		testDiff(t, "Comments", expected[i].Comments, actual[i].Comments)
		testDiff(t, "ReservedTags", expected[i].ReservedTags, actual[i].ReservedTags)
		testDiff(t, "ReservedNames", expected[i].ReservedNames, actual[i].ReservedNames)
	}
}

func TestFields(t *testing.T) {
	var (
		expected = []Field{
			Field{
				id:     fieldID(".testv2.Message:string_field"),
				parent: messageID(".testv2.Message"),
				Name:   "string_field",
				Meta: meta.FieldMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Generator:  "email",
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.FieldOptions{
					Deprecated: boolPointer(true),
				},
				Comments: Comments{
					Leading: " string_field comment\n",
				},
				Number:       1,
				Label:        descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:         descriptor.FieldDescriptorProto_TYPE_STRING,
				DefaultValue: "",
				JSONName:     "stringField",
			},
			Field{
				id:     fieldID(".testv2.Message:repeated_string_field"),
				parent: messageID(".testv2.Message"),
				Name:   "repeated_string_field",
				Comments: Comments{
					Leading: " repeated_string_field comment\n",
				},
				Number:   2,
				Label:    descriptor.FieldDescriptorProto_LABEL_REPEATED,
				Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
				JSONName: "repeatedStringField",
			},
			Field{
				id:       fieldID(".testv2.Message:enum_field"),
				parent:   messageID(".testv2.Message"),
				typeEnum: enumID(".testv2.Enum"),
				Name:     "enum_field",
				Comments: Comments{
					Leading: " enum_field comment\n",
				},
				Number:   3,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_ENUM,
				JSONName: "enumField",
			},
			Field{
				id:          fieldID(".testv2.Message:other_message_field"),
				parent:      messageID(".testv2.Message"),
				typeMessage: messageID(".testv2.OtherMessage"),
				Name:        "other_message_field",
				Comments: Comments{
					Leading: " other_message_field comment\n",
				},
				Number:   4,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE,
				JSONName: "otherMessageField",
			},
			Field{
				id:     fieldID(".testv2.Message:bool_field"),
				parent: messageID(".testv2.Message"),
				oneof:  oneofID(".testv2.Message:oneof_field"),
				Name:   "bool_field",
				Comments: Comments{
					Leading: " bool_field comment\n",
				},
				Number:   5,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_BOOL,
				JSONName: "boolField",
			},
			Field{
				id:       fieldID(".testv2.Message:embedded_enum_field"),
				parent:   messageID(".testv2.Message"),
				typeEnum: enumID(".testv2.Message.EmbeddedEnum"),
				Name:     "embedded_enum_field",
				Comments: Comments{
					Leading: " embedded_enum_field comment\n",
				},
				Number:   6,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_ENUM,
				JSONName: "embeddedEnumField",
			},
			Field{
				id:          fieldID(".testv2.Message:embedded_message_field"),
				parent:      messageID(".testv2.Message"),
				typeMessage: messageID(".testv2.Message.EmbeddedMessage"),
				Name:        "embedded_message_field",
				Comments: Comments{
					Leading: " embedded_message_field comment\n",
				},
				Number:   7,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE,
				JSONName: "embeddedMessageField",
			},
			Field{
				id:     fieldID(".testv2.Message.EmbeddedMessage:string_field"),
				parent: messageID(".testv2.Message.EmbeddedMessage"),
				Name:   "string_field",
				Comments: Comments{
					Leading: " string_field comment\n",
				},
				Number:   1,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
				JSONName: "stringField",
			},
			Field{
				id:     fieldID(".testv2.Message.EmbeddedMessage:uint32_field"),
				parent: messageID(".testv2.Message.EmbeddedMessage"),
				Name:   "uint32_field",
				Comments: Comments{
					Trailing: " uint32_field trailing comment\n",
				},
				Number:   2,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_UINT32,
				JSONName: "uint32Field",
			},
			Field{
				id:     fieldID(".testv3.Message:string_field"),
				parent: messageID(".testv3.Message"),
				Name:   "string_field",
				Meta: meta.FieldMetadata{
					Visibility: meta.Visibility_PRIVATE,
					Generator:  "email",
					Tags:       []string{"tag1", "tag2"},
					Extra:      map[string]string{"k": "v"},
				},
				Options: descriptor.FieldOptions{
					Deprecated: boolPointer(true),
				},
				Comments: Comments{
					Leading: " string_field comment\n",
				},
				Number:       1,
				Label:        descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:         descriptor.FieldDescriptorProto_TYPE_STRING,
				DefaultValue: "",
				JSONName:     "stringField",
			},
			Field{
				id:     fieldID(".testv3.Message:repeated_string_field"),
				parent: messageID(".testv3.Message"),
				Name:   "repeated_string_field",
				Comments: Comments{
					Leading: " repeated_string_field comment\n",
				},
				Number:   2,
				Label:    descriptor.FieldDescriptorProto_LABEL_REPEATED,
				Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
				JSONName: "repeatedStringField",
			},
			Field{
				id:       fieldID(".testv3.Message:enum_field"),
				parent:   messageID(".testv3.Message"),
				typeEnum: enumID(".testv3.Enum"),
				Name:     "enum_field",
				Comments: Comments{
					Leading: " enum_field comment\n",
				},
				Number:   3,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_ENUM,
				JSONName: "enumField",
			},
			Field{
				id:          fieldID(".testv3.Message:other_message_field"),
				parent:      messageID(".testv3.Message"),
				typeMessage: messageID(".testv3.OtherMessage"),
				Name:        "other_message_field",
				Comments: Comments{
					Leading: " other_message_field comment\n",
				},
				Number:   4,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE,
				JSONName: "otherMessageField",
			},
			Field{
				id:     fieldID(".testv3.Message:bool_field"),
				parent: messageID(".testv3.Message"),
				oneof:  oneofID(".testv3.Message:oneof_field"),
				Name:   "bool_field",
				Comments: Comments{
					Leading: " bool_field comment\n",
				},
				Number:   5,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_BOOL,
				JSONName: "boolField",
			},
			Field{
				id:       fieldID(".testv3.Message:embedded_enum_field"),
				parent:   messageID(".testv3.Message"),
				typeEnum: enumID(".testv3.Message.EmbeddedEnum"),
				Name:     "embedded_enum_field",
				Comments: Comments{
					Leading: " embedded_enum_field comment\n",
				},
				Number:   6,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_ENUM,
				JSONName: "embeddedEnumField",
			},
			Field{
				id:          fieldID(".testv3.Message:embedded_message_field"),
				parent:      messageID(".testv3.Message"),
				typeMessage: messageID(".testv3.Message.EmbeddedMessage"),
				Name:        "embedded_message_field",
				Comments: Comments{
					Leading: " embedded_message_field comment\n",
				},
				Number:   7,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE,
				JSONName: "embeddedMessageField",
			},
			Field{
				id:     fieldID(".testv3.Message.EmbeddedMessage:string_field"),
				parent: messageID(".testv3.Message.EmbeddedMessage"),
				Name:   "string_field",
				Comments: Comments{
					Leading: " string_field comment\n",
				},
				Number:   1,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
				JSONName: "stringField",
			},
			Field{
				id:     fieldID(".testv3.Message.EmbeddedMessage:uint32_field"),
				parent: messageID(".testv3.Message.EmbeddedMessage"),
				Name:   "uint32_field",
				Comments: Comments{
					Trailing: " uint32_field trailing comment\n",
				},
				Number:   2,
				Label:    descriptor.FieldDescriptorProto_LABEL_OPTIONAL,
				Type:     descriptor.FieldDescriptorProto_TYPE_UINT32,
				JSONName: "uint32Field",
			},
		}
		actual = New(&request).ToGenerate().Fields()
	)

	for i := 0; i < len(actual); i++ {
		testDiff(t, "id", expected[i].id, actual[i].id)
		testDiff(t, "parent", expected[i].parent, actual[i].parent)
		testDiff(t, "oneof", expected[i].oneof, actual[i].oneof)
		testDiff(t, "typeMessage", expected[i].typeMessage, actual[i].typeMessage)
		testDiff(t, "typeEnum", expected[i].typeEnum, actual[i].typeEnum)
		testDiff(t, "Name", expected[i].Name, actual[i].Name)
		testDiff(t, "Meta", expected[i].Meta, actual[i].Meta)
		testDiff(t, "Options.Deprecated", expected[i].Options.Deprecated, actual[i].Options.Deprecated)
		testDiff(t, "Comments", expected[i].Comments, actual[i].Comments)
		testDiff(t, "Number", expected[i].Number, actual[i].Number)
		testDiff(t, "Label", expected[i].Label, actual[i].Label)
		testDiff(t, "Type", expected[i].Type, actual[i].Type)
		testDiff(t, "DefaultValue", expected[i].DefaultValue, actual[i].DefaultValue)
		testDiff(t, "JSONName", expected[i].JSONName, actual[i].JSONName)
	}
}

func testDiff(t *testing.T, sbj string, expected, actual interface{}) {
	diffs := pretty.Diff(expected, actual)
	for _, diff := range diffs {
		log.Printf("Mismatch: %s", diff)
	}
	if len(diffs) > 0 {
		t.Fatalf("%s mismatches", sbj)
	}
}
