package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// Funcs is the template.FuncMap used for template execution
var Funcs = template.FuncMap{
	"exec":       Exec,
	"gofmt":      GoFmt,
	"uppercamel": UpperCamel,
	"upper":      Upper,
	"lowercamel": LowerCamel,
	"lower":      Lower,
	"replace":    Replace,
	"join":       strings.Join,
	"split":      strings.Split,
	"dict":       Dict,
	"merge":      Merge,
	"base64":     base64.StdEncoding.EncodeToString,
}

// Exec executes the named template, returning its output as a string
func Exec(name string, data interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, name, data)
	return buf.String(), err
}

// GoFmt applies gofmt to the string
func GoFmt(s string) (string, error) {
	b, err := format.Source([]byte(s))
	if err != nil {
		return "", fmt.Errorf("%s\n%s", err, s)
	}
	return string(b), nil
}

// UpperCamel converts a snake_case string into CamelCase
func UpperCamel(s string) string {
	return generator.CamelCase(s)
}

// Upper converts a string to uppercase
func Upper(s string) string {
	return strings.ToUpper(s)
}

// LowerCamel converts a snake_case string into camelCase
func LowerCamel(s string) string {
	if len(s) == 0 {
		return s
	}
	b := []byte(generator.CamelCase(s))
	b[0] = bytes.ToLower(b[0:1])[0]
	return string(b)
}

// Lower converts a string to lowercase
func Lower(s string) string {
	return strings.ToLower(s)
}

// Replace wraps strings.Replace for better use with template pipes
// Example: "foo" | replace "f" "p" | replace "p" "m"
func Replace(old, new, s string) string {
	return strings.Replace(s, old, new, -1)
}

// Dict converts a set of name/value pairs into a map
//
// Example:
//
//   {{ $map := dict "Foo" 1 "Bar" 2 }}
//   {{ $map.Foo }} -> 1
//   {{ $map.Bar }} -> 2
//
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// Merge returns a copy of the given map with the given set of name/value pairs
// assigned
//
// Example:
//
//   {{ $map := dict "Foo" 1 "Bar" 2 }}
//   {{ $map.Foo }} -> 1
//   {{ $map.Bar }} -> 2
//   {{ (merge $map "Foo" 3).Foo }} -> 3
//   {{ (merge $map "Foo" 3).Bar }} -> 2
//   {{ $map.Foo }} -> 1
//
func Merge(source map[string]interface{}, values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(source)+(len(values)/2))
	for k, v := range source {
		dict[k] = v
	}

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}
