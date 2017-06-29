package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}

	if len(data) == 0 {
		log.Fatalf("Empty input")
	}

	var response plugin.CodeGeneratorResponse
	if err := proto.Unmarshal(data, &response); err != nil {
		log.Fatalf("Error unmarshaling input: %s", err)
	}

	if response.Error != nil {
		log.Fatalf("Error response: %s", *response.Error)
	}

	for _, file := range response.File {
		if file.Name != nil {
			log.Printf("Name: %s\n", *file.Name)
		}
		if file.InsertionPoint != nil {
			log.Printf("InsertionPoint: %s\n", *file.InsertionPoint)
		}
		if file.Content != nil {
			log.Printf("Content: \n%s\n", *file.Content)
		}
	}
	if len(response.File) == 0 {
		log.Printf("(File empty)")
	}
}
