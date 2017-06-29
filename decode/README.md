# dump prints `plugin.CodeGeneratorResponse` to STDOUT as JSON

Useful for inspecting the output of protoc plugins

Example:

```sh
protoc --dump_out=. example.proto
cat dump.pb | protoc-gen-go | go run main.go
```
