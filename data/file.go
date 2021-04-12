package data

import (
	"path"
	"sort"
	"strings"

	"github.com/kerinin/protoc-gen-template/meta"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type fileID string

// FileSlice is a slice of files
type FileSlice []File

// Visible returns the visible values in the slice
func (s FileSlice) Visible() FileSlice {
	outputs := make([]File, 0, len(s))
	for _, f := range s {
		if f.IsVisible() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// NotDeprecated returns the non-deprecated values in the slice
func (s FileSlice) NotDeprecated() FileSlice {
	outputs := make([]File, 0, len(s))
	for _, f := range s {
		if !f.IsDeprecated() {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// ToGenerate returns the values in the slice to be generated
func (s FileSlice) ToGenerate() FileSlice {
	outputs := make([]File, 0, len(s))
	for _, f := range s {
		if f.Generate {
			outputs = append(outputs, f)
		}
	}
	return outputs
}

// File describes a protobuf source file
type File struct {
	idx            int
	id             fileID
	data           *Data
	messages       []messageID
	enums          []enumID
	services       []serviceID
	sourceCodeInfo map[string]*descriptor.SourceCodeInfo_Location

	Package      string
	Name         string
	Comments     Comments
	Meta         meta.FileMetadata      // Custom metadata extensions defined for protoc-gen-template
	Options      descriptor.FileOptions // Globally-defined file metadata
	Generate     bool                   // Generate is true if the file is included in FileToGenerate
	Dependencies []string               // Dependencies lists the names of files imported by this file.
	Syntax       string                 // Syntax is of the proto file - proto2/proto3
}

func (f File) String() string {
	return string(f.id)
}

// IsVisible returns true if the visibility metadata of the file is `PUBLIC`
func (f File) IsVisible() bool {
	return f.Meta.Visibility == meta.Visibility_PUBLIC
}

// IsDeprecated returns true if the file's deprecation option is true
func (f File) IsDeprecated() bool {
	return f.Options.Deprecated != nil && *f.Options.Deprecated == true
}

// GoPackageName returns a package name for importing the file:
//
// 1. If the file's GoPackage option is set and contains the `;` character,
//    returns any content after the character
// 2. If the file's GoPackage options is set and doesn't contain the `;`
//    character, returns the basename portion of the GoPackage value
// 3. Returns the basename portion of file's Name value
func (f File) GoPackageName() string {
	if f.Options.GoPackage == nil || *f.Options.GoPackage == "" {
		return path.Base(f.Name)
	}
	parts := strings.Split(*f.Options.GoPackage, `;`)
	if len(parts) > 1 {
		return parts[1]
	}
	return path.Base(*f.Options.GoPackage)
}

// GoPackageImport returns the go import path for the file's package
//
// 1. If the file's GoPackage option is set and contains the `;` character,
//    returns any content before the character
// 2. If the file's GoPackage options is set and doesn't contain the `;`
//    character, returns the GoPackage value
// 3. Returns the directory portion of the file's Name value
// 4. Returns the file's Name value
func (f File) GoPackageImport() string {
	if f.Options.GoPackage == nil || *f.Options.GoPackage == "" {
		if path.Dir(f.Name) == "." {
			return f.Name
		}
		return path.Dir(f.Name)
	}
	return strings.Split(*f.Options.GoPackage, `;`)[0]
}

// Messages returns a slice of the file's messages
func (f File) Messages() MessageSlice {
	vs := make([]Message, 0, len(f.messages))
	for _, v := range f.messages {
		vs = append(vs, *f.data.messages[v])
	}
	sort.Sort(sortedMessagesByIndex(vs))
	return vs
}

// Enums returns a slice of the file's enums
func (f File) Enums() EnumSlice {
	vs := make([]Enum, 0, len(f.enums))
	for _, v := range f.enums {
		vs = append(vs, *f.data.enums[v])
	}
	return vs
}

// Services returns a slice of the file's services
func (f File) Services() ServiceSlice {
	vs := make([]Service, 0, len(f.services))
	for _, v := range f.services {
		vs = append(vs, *f.data.services[v])
	}
	return vs
}

func newFileMetadata(in *descriptor.FileOptions) (out meta.FileMetadata) {
	defer func() {
		// NOTE: There's a bug in `proto` that causes panics when calling
		// `GetExtension`, `HasExtension`, etc in some cases when there isn't
		// a defined extension.  This recovers from the panic and allows a
		// `nil` return.
		_ = recover()
	}()

	ext := proto.GetExtension(in, meta.E_FileMeta)
	return *ext.(*meta.FileMetadata)
}

func derefFileOptions(o *descriptor.FileOptions) (_ descriptor.FileOptions) {
	if o == nil {
		return
	}
	return *o
}

type sortedFilesByIndex []File

func (s sortedFilesByIndex) Len() int           { return len(s) }
func (s sortedFilesByIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedFilesByIndex) Less(i, j int) bool { return s[i].idx < s[j].idx }
