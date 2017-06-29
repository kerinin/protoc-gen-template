package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ReturnPath/protoc-gen-template/data"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

var (
	debugTemplate   = flag.String("template", "", "Template file (for debugging purposes)")
	defaultTemplate = "."
	tmpl            *template.Template
)

type fileInfo struct {
	inPath       string
	outPath      string
	templateName string
}

func generateFiles(req *plugin.CodeGeneratorRequest) ([]*plugin.CodeGeneratorResponse_File, error) {

	if len(req.FileToGenerate) == 0 {
		return nil, errors.New("no files to generate")
	}

	if *debugTemplate != "" {
		req.Parameter = debugTemplate
	}
	if req.Parameter == nil || *req.Parameter == "" {
		req.Parameter = &defaultTemplate
	}

	var err error

	// NOTE: `tmpl` is a global varible so it can be accessed from inside
	// functions passed to the functionmap.  Specifically, `exec` needs access
	// to the template map to be able to render & capture template output.  In
	// order for `exec` to be present in the search path it needs to be defined
	// before we parse templates, but in order for it to actually exec templates
	// it needs access to all the compiled output.  It's an ugly solution but
	// is the only one I can find ATM.
	tmpl = template.New("").Funcs(Funcs)

	var (
		templateFiles = []fileInfo{}
		copyFiles     = []fileInfo{}
	)
	err = filepath.Walk(*req.Parameter, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "walking at path %s", filename)
		}
		if info.IsDir() {
			return nil
		}

		// The path relative to the compile output directory
		outPath, err := filepath.Rel(*req.Parameter, filename)
		if err != nil {
			return errors.Wrap(err, "building relative path")
		}

		switch {
		case strings.HasSuffix(info.Name(), `.associated.tmpl`):
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return errors.Wrapf(err, "reading file %s", filename)
			}

			templateName := strings.TrimSuffix(filepath.ToSlash(outPath), `.associated.tmpl`)

			tmpl = tmpl.New(templateName)
			tmpl, err = tmpl.Parse(string(b))
			if err != nil {
				return errors.Wrapf(err, "parsing template %s", filename)
			}

		case strings.HasSuffix(info.Name(), `.tmpl`):
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return errors.Wrapf(err, "reading file %s", filename)
			}

			info := fileInfo{
				inPath:       filename,
				outPath:      strings.TrimSuffix(outPath, `.tmpl`),
				templateName: strings.TrimSuffix(filepath.ToSlash(outPath), `.tmpl`),
			}

			tmpl = tmpl.New(info.templateName)
			tmpl, err = tmpl.Parse(string(b))
			if err != nil {
				return errors.Wrapf(err, "parsing template %s", filename)
			}
			templateFiles = append(templateFiles, info)

		default:
			copyFiles = append(copyFiles, fileInfo{
				inPath:  filename,
				outPath: outPath,
			})
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "walking input")
	}

	files := make([]*plugin.CodeGeneratorResponse_File, 0, len(templateFiles)+len(copyFiles))

	for _, f := range templateFiles {
		file, err := generateFile(f, data.New(req))
		if err != nil {
			return nil, errors.Wrapf(err, "generating file %s", f.inPath)
		}
		files = append(files, file)
	}

	for _, f := range copyFiles {
		file, err := copyFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "copying file %s", f.inPath)
		}
		files = append(files, file)
	}

	return files, nil
}

func generateFile(f fileInfo, d *data.Data) (*plugin.CodeGeneratorResponse_File, error) {
	buffer := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buffer, f.templateName, d)
	if err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	content := buffer.String()
	return &plugin.CodeGeneratorResponse_File{
		Name:    &f.outPath,
		Content: &content,
	}, nil
}

func copyFile(f fileInfo) (*plugin.CodeGeneratorResponse_File, error) {
	b, err := ioutil.ReadFile(f.inPath)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", f.inPath)
	}
	s := string(b)

	return &plugin.CodeGeneratorResponse_File{
		Name:    &f.outPath,
		Content: &s,
	}, nil
}

func fail(err error) {
	log.Print("protoc-gen-render-template: error:", err)
	os.Exit(1)
}

func respondFail(err error) {
	msg := err.Error()
	res := &plugin.CodeGeneratorResponse{
		Error: &msg,
	}

	data, err := proto.Marshal(res)
	if err != nil {
		fail(errors.Wrap(err, "marshaling output"))
	}

	_, _ = os.Stdout.Write(data)
	os.Exit(0)
}

func main() {
	flag.Parse()

	var (
		req plugin.CodeGeneratorRequest
		res plugin.CodeGeneratorResponse
	)

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fail(errors.Wrap(err, "reading STDIN"))
	}

	if err := proto.Unmarshal(data, &req); err != nil {
		fail(errors.Wrap(err, "unmarshaling STDIN"))
	}

	files, err := generateFiles(&req)

	if err != nil {
		respondFail(errors.Wrap(err, "generating files"))
	}

	res.File = files

	data, err = proto.Marshal(&res)
	if err != nil {
		fail(errors.Wrap(err, "marshaling output"))
	}

	_, _ = os.Stdout.Write(data)
}
