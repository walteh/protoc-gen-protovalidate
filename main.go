package main

import (
	"embed"

	"github.com/bufbuild/protoplugin"
	protocgenprotovalidate "github.com/walteh/protoc-gen-protovalidate/cmd/protoc-gen-protovalidate"
)

//go:generate go run ./generator

//go:embed gen/protovalidate-*-latest*
var localFiles embed.FS

func main() {
	handler := &protocgenprotovalidate.Handler{
		Cache: localFiles,
	}
	protoplugin.Main(protoplugin.HandlerFunc(handler.Handle))
}
