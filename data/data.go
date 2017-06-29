package data

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Data is type assigned to "dot" when evaluating templates
type Data struct {
	enums      map[enumID]*Enum
	enumValues map[enumValueID]*EnumValue
	fields     map[fieldID]*Field
	files      map[fileID]*File
	messages   map[messageID]*Message
	methods    map[methodID]*Method
	oneofs     map[oneofID]*Oneof
	services   map[serviceID]*Service

	filesToGenerate map[string]bool
	fileCount       int
	msgCount        int
	fieldCount      int
}

// New returns a new Data describing the code generator request
func New(req *plugin.CodeGeneratorRequest) *Data {
	data := &Data{
		enums:           map[enumID]*Enum{},
		enumValues:      map[enumValueID]*EnumValue{},
		fields:          map[fieldID]*Field{},
		files:           make(map[fileID]*File, len(req.ProtoFile)),
		messages:        map[messageID]*Message{},
		methods:         map[methodID]*Method{},
		oneofs:          map[oneofID]*Oneof{},
		services:        map[serviceID]*Service{},
		filesToGenerate: make(map[string]bool, len(req.FileToGenerate)),
	}

	// Build files to generate index
	for _, file := range req.FileToGenerate {
		data.filesToGenerate[file] = true
	}

	// Merge files & their contents
	for _, file := range req.ProtoFile {
		data.mergeFile(file)
	}

	return data
}

func (d *Data) mergeEnum(f fileID, m messageID, desc *descriptor.EnumDescriptorProto, path string) enumID {
	enum := &Enum{
		data:     d,
		file:     f,
		parent:   m,
		values:   make([]enumValueID, 0, len(desc.Value)),
		Name:     *desc.Name,
		Meta:     newEnumMetadata(desc.Options),
		Options:  derefEnumOptions(desc.Options),
		Comments: d.comments(f, path),
	}

	if m == "" {
		enum.id = enumID(fmt.Sprintf(".%s.%s", d.files[f].Package, *desc.Name))
	} else {
		// Nested enums use their parent's name as a prefix
		enum.id = enumID(fmt.Sprintf("%s.%s", m, *desc.Name))
	}

	for i, desc := range desc.Value {
		// Value is field 2 in EnumDescriptorProto
		p := fmt.Sprintf("%s,2,%d", path, i)
		id := d.mergeEnumValue(f, enum.id, desc, p)
		enum.values = append(enum.values, id)
	}

	d.enums[enum.id] = enum
	return enum.id
}

func (d *Data) mergeEnumValue(f fileID, e enumID, desc *descriptor.EnumValueDescriptorProto, path string) enumValueID {
	enumValue := &EnumValue{
		id:       enumValueID(fmt.Sprintf("%s:%s", e, *desc.Name)),
		data:     d,
		parent:   e,
		Name:     *desc.Name,
		Meta:     newEnumValueMetadata(desc.Options),
		Options:  derefEnumValueOptions(desc.Options),
		Comments: d.comments(f, path),
		Number:   *desc.Number,
	}

	d.enumValues[enumValue.id] = enumValue
	return enumValue.id
}

func (d *Data) mergeField(f fileID, m *Message, desc *descriptor.FieldDescriptorProto, path string) fieldID {
	field := &Field{
		idx:          d.fieldCount,
		id:           fieldID(fmt.Sprintf("%s:%s", m, *desc.Name)),
		data:         d,
		parent:       m.id,
		Type:         *desc.Type,
		Name:         *desc.Name,
		Meta:         newFieldMetadata(desc.Options),
		Options:      derefFieldOptions(desc.Options),
		Comments:     d.comments(f, path),
		Number:       *desc.Number,
		Label:        *desc.Label,
		DefaultValue: toString(desc.DefaultValue, ""),
		JSONName:     *desc.JsonName,
	}
	d.fieldCount++
	d.fields[field.id] = field

	// Associate oneof fields with their associated Oneof
	if desc.OneofIndex != nil {
		id := m.oneofs[*desc.OneofIndex]
		field.oneof = id
		d.oneofs[id].fields = append(d.oneofs[id].fields, field.id)
	}

	// NOTE: This may fail to locate types - from the docs:
	//
	//   For message and enum types, this is the name of the type.  if the name
	//   starts with a '.', it is fully-qualified.  otherwise, c++-like scoping
	//   rules are used to find the type (i.e. first the nested types within this
	//   message are searched, then within the parent, on up to the root
	//   namespace).
	//
	if field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		field.typeMessage = messageID(*desc.TypeName)
	}
	if field.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		field.typeEnum = enumID(*desc.TypeName)
	}

	return field.id
}

func (d *Data) mergeFile(desc *descriptor.FileDescriptorProto) {
	file := &File{
		idx:            d.fileCount,
		id:             fileID(fmt.Sprintf(".%s:%s", *desc.Package, *desc.Name)),
		data:           d,
		messages:       make([]messageID, 0, len(desc.MessageType)),
		enums:          make([]enumID, 0, len(desc.EnumType)),
		services:       make([]serviceID, 0, len(desc.Service)),
		sourceCodeInfo: make(map[string]*descriptor.SourceCodeInfo_Location, len(desc.SourceCodeInfo.Location)),
		Name:           *desc.Name,
		Meta:           newFileMetadata(desc.Options),
		Options:        derefFileOptions(desc.Options),
		Package:        *desc.Package,
		Generate:       d.filesToGenerate[*desc.Name],
		Dependencies:   desc.Dependency,
		Syntax:         toString(desc.Syntax, "proto2"),
	}
	d.fileCount++
	d.files[file.id] = file

	// Build source code index
	for _, l := range desc.SourceCodeInfo.Location {
		pathParts := make([]string, 0, len(l.Path))
		for _, part := range l.Path {
			pathParts = append(pathParts, fmt.Sprintf("%d", part))
		}
		file.sourceCodeInfo[strings.Join(pathParts, ",")] = l
	}

	// Package is field 2 in FileDescriptorProto
	file.Comments = d.comments(file.id, "2")

	for i, dsc := range desc.MessageType {
		// MessageType is field 4 in FileDescriptorProto
		p := fmt.Sprintf("4,%d", i)
		id := d.mergeMessage(file.id, "", dsc, p)
		file.messages = append(file.messages, id)
	}

	for i, dsc := range desc.EnumType {
		// EnumType is field 5 in FileDescriptorProto
		p := fmt.Sprintf("5,%d", i)
		id := d.mergeEnum(file.id, "", dsc, p)
		file.enums = append(file.enums, id)
	}

	for i, dsc := range desc.Service {
		// Service is field 6 in FileDescriptorProto
		p := fmt.Sprintf("6,%d", i)
		id := d.mergeService(file.id, dsc, p)
		file.services = append(file.services, id)
	}
}

func (d *Data) mergeMessage(f fileID, m messageID, desc *descriptor.DescriptorProto, path string) messageID {
	message := &Message{
		idx:           d.msgCount,
		data:          d,
		file:          f,
		parent:        m,
		fields:        make([]fieldID, 0, len(desc.Field)),
		messages:      make([]messageID, 0, len(desc.NestedType)),
		enums:         make([]enumID, 0, len(desc.EnumType)),
		oneofs:        make([]oneofID, 0, len(desc.OneofDecl)),
		Name:          *desc.Name,
		Meta:          newMessageMetadata(desc.Options),
		Options:       derefMessageOptions(desc.Options),
		Comments:      d.comments(f, path),
		ReservedTags:  make([]descriptor.DescriptorProto_ReservedRange, 0, len(desc.ReservedRange)),
		ReservedNames: desc.ReservedName,
	}
	d.msgCount++

	if m == "" {
		message.id = messageID(fmt.Sprintf(".%s.%s", d.files[f].Package, message.Name))
	} else {
		// Nested messages use their parent's name as a prefix
		message.id = messageID(fmt.Sprintf("%s.%s", m, message.Name))
	}
	d.messages[message.id] = message

	// Derefernce pointers
	for _, tag := range desc.ReservedRange {
		message.ReservedTags = append(message.ReservedTags, *tag)
	}

	for i, dsc := range desc.OneofDecl {
		// OneofDecl is field 8 in DescriptorProto
		p := fmt.Sprintf("%s,8,%d", path, i)
		id := d.mergeOneof(f, message.id, dsc, p)
		message.oneofs = append(message.oneofs, id)
	}

	for i, dsc := range desc.Field {
		// Field is field 2 in DescriptorProto
		p := fmt.Sprintf("%s,2,%d", path, i)
		id := d.mergeField(f, message, dsc, p)
		message.fields = append(message.fields, id)
	}

	for i, dsc := range desc.NestedType {
		// NestedType is field 3 in DescriptorProto
		p := fmt.Sprintf("%s,3,%d", path, i)
		id := d.mergeMessage(f, message.id, dsc, p)
		message.messages = append(message.messages, id)
	}

	for i, dsc := range desc.EnumType {
		// EnumType is field 4 in DescriptorProto
		p := fmt.Sprintf("%s,4,%d", path, i)
		id := d.mergeEnum(f, message.id, dsc, p)
		message.enums = append(message.enums, id)
	}

	return message.id
}

func (d *Data) mergeMethod(f fileID, s serviceID, desc *descriptor.MethodDescriptorProto, path string) methodID {
	method := &Method{
		id:              methodID(fmt.Sprintf("%s:%s", s, *desc.Name)),
		data:            d,
		parent:          s,
		inputType:       messageID(*desc.InputType),
		outputType:      messageID(*desc.OutputType),
		Name:            *desc.Name,
		Meta:            newMethodMetadata(desc.Options),
		Options:         derefMethodOptions(desc.Options),
		Comments:        d.comments(f, path),
		ClientStreaming: derefBool(desc.ClientStreaming),
		ServerStreaming: derefBool(desc.ServerStreaming),
	}

	// NOTE: From the proto comments:
	//
	//   Input and output type names.  These are resolved in the same way as
	//   FieldDescriptorProto.type_name, but must refer to a message type.
	//
	//   For message and enum types, this is the name of the type.  If the name
	//   starts with a '.', it is fully-qualified.  Otherwise, C++-like scoping
	//   rules are used to find the type (i.e. first the nested types within this
	//   message are searched, then within the parent, on up to the root
	//   namespace).
	//
	if _, found := d.messages[method.inputType]; !found {
		keys := make([]messageID, 0, len(d.messages))
		for k := range d.messages {
			keys = append(keys, k)
		}
		panic(fmt.Sprintf("No key %s in %v", method.inputType, keys))
	}
	if _, found := d.messages[method.outputType]; !found {
		keys := make([]messageID, 0, len(d.messages))
		for k := range d.messages {
			keys = append(keys, k)
		}
		panic(fmt.Sprintf("No key %s in %v", method.outputType, keys))
	}

	d.methods[method.id] = method
	return method.id
}

func (d *Data) mergeOneof(f fileID, m messageID, desc *descriptor.OneofDescriptorProto, path string) oneofID {
	oneof := &Oneof{
		id:       oneofID(fmt.Sprintf("%s:%s", m, *desc.Name)),
		data:     d,
		parent:   m,
		Name:     *desc.Name,
		Options:  derefOneofOptions(desc.Options),
		Comments: d.comments(f, path),
	}

	d.oneofs[oneof.id] = oneof
	return oneof.id
}

func (d *Data) mergeService(f fileID, desc *descriptor.ServiceDescriptorProto, path string) serviceID {
	service := &Service{
		id:       serviceID(fmt.Sprintf(".%s.%s", d.files[f].Package, *desc.Name)),
		data:     d,
		file:     f,
		methods:  make([]methodID, 0, len(desc.Method)),
		Name:     *desc.Name,
		Meta:     newServiceMetadata(desc.Options),
		Options:  derefServiceOptions(desc.Options),
		Comments: d.comments(f, path),
	}

	for i, dsc := range desc.Method {
		// Method is field 2 in ServiceDescriptorProto
		p := fmt.Sprintf("%s,2,%d", path, i)
		id := d.mergeMethod(f, service.id, dsc, p)
		service.methods = append(service.methods, id)
	}

	d.services[service.id] = service
	return service.id
}

func (d *Data) comments(f fileID, path string) Comments {
	location, found := d.files[f].sourceCodeInfo[path]
	if !found {
		return Comments{}
	}

	return Comments{
		Leading:         toString(location.LeadingComments, ""),
		Trailing:        toString(location.TrailingComments, ""),
		LeadingDetached: location.LeadingDetachedComments,
	}
}

// Enums returns a slice of defined enums
func (d *Data) Enums() EnumSlice {
	vs := make([]Enum, 0, len(d.enums))
	for _, v := range d.enums {
		vs = append(vs, *v)
	}
	return vs
}

// EnumValues returns a slice of defined enum values
func (d *Data) EnumValues() EnumValueSlice {
	vs := make([]EnumValue, 0, len(d.enumValues))
	for _, v := range d.enumValues {
		vs = append(vs, *v)
	}
	return vs
}

// Fields returns a slice of defined fields
func (d *Data) Fields() FieldSlice {
	vs := make([]Field, 0, len(d.fields))
	for _, v := range d.fields {
		if v.Type == descriptor.FieldDescriptorProto_TYPE_GROUP {
			continue
		}
		vs = append(vs, *v)
	}
	sort.Sort(sortedFieldsByIndex(vs))
	return vs
}

// Files returns a slice of defined files
func (d *Data) Files() FileSlice {
	vs := make([]File, 0, len(d.files))
	for _, v := range d.files {
		vs = append(vs, *v)
	}
	sort.Sort(sortedFilesByIndex(vs))
	return vs
}

// Messages returns a slice of defined messages
func (d *Data) Messages() MessageSlice {
	vs := make([]Message, 0, len(d.messages))
	for _, v := range d.messages {
		vs = append(vs, *v)
	}
	sort.Sort(sortedMessagesByIndex(vs))
	return vs
}

// Methods returns a slice of defined methods
func (d *Data) Methods() MethodSlice {
	vs := make([]Method, 0, len(d.methods))
	for _, v := range d.methods {
		vs = append(vs, *v)
	}
	return vs
}

// Oneofs returns a slice of defined oneofs
func (d *Data) Oneofs() OneofSlice {
	vs := make([]Oneof, 0, len(d.oneofs))
	for _, v := range d.oneofs {
		vs = append(vs, *v)
	}
	return vs
}

// Services returns a slice of defined services
func (d *Data) Services() ServiceSlice {
	vs := make([]Service, 0, len(d.services))
	for _, v := range d.services {
		vs = append(vs, *v)
	}
	return vs
}

// PackagesToGenerate returns a slice containing all package names used in files to generate
func (d *Data) PackagesToGenerate() []string {
	packageMap := make(map[string]struct{}, len(d.files))

	for _, file := range d.files {
		if file.Generate {
			packageMap[file.Package] = struct{}{}
		}
	}

	packages := make([]string, 0, len(packageMap))
	for pkg := range packageMap {
		packages = append(packages, pkg)
	}

	return packages
}

func toString(s *string, d string) string {
	if s == nil {
		return d
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return true
}
