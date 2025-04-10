package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bufbuild/protoplugin"

	"google.golang.org/protobuf/types/descriptorpb"
)

//go:generate go run ./generator

//go:embed  templates
var root embed.FS

func AllGoFiles() (embed.FS, error) {
	return root, nil
}

func main() {
	protoplugin.Main(protoplugin.HandlerFunc(handle))
}

func handle(
	ctx context.Context,
	e protoplugin.PluginEnv,
	responseWriter protoplugin.ResponseWriter,
	request protoplugin.Request,
) error {

	responseWriter.SetFeatureProto3Optional()
	responseWriter.SetFeatureSupportsEditions(descriptorpb.Edition_EDITION_2023, descriptorpb.Edition_EDITION_2024)

	fileDescriptors, err := request.FileDescriptorsToGenerate()
	if err != nil {
		return err
	}

	params, err := parseParams(request.Parameter())
	if err != nil {
		return err
	}

	for _, fileDescriptor := range fileDescriptors {

		if !strings.HasSuffix(fileDescriptor.Path(), params.GetBufValidateFile()) {
			continue
		}

		desc, ok := fileDescriptor.Options().(*descriptorpb.FileOptions)
		if !ok {
			return fmt.Errorf("file descriptor is not a FileDescriptorProto")
		}

		gpkg := desc.GetGoPackage()

		// fmt.Fprintf(os.Stderr, "fileDescriptor.Path(): %+v\n", fileDescriptor.Path())
		// fmt.Fprintf(os.Stderr, "fileDescriptor.Name(): %+v\n", fileDescriptor.Name())
		// fmt.Fprintf(os.Stderr, "fileDescriptor.Package(): %+v\n", fileDescriptor.Package())
		// fmt.Fprintf(os.Stderr, "fileDescriptor.FullName(): %+v\n", fileDescriptor.FullName())

		rootPath := strings.TrimSuffix(fileDescriptor.Path(), "validate.proto")

		// fmt.Fprintf(os.Stderr, "gpkg: %+v\n", gpkg)

		tmpl := template.New("root")

		tmpl.Delims("[[[[[[[[", "]]]]]]]]")

		tmpl, err = tmpl.ParseFS(root, "templates/*.tmpl")
		if err != nil {
			return err
		}

		for _, file := range tmpl.Templates() {
			if file.Name() == "root" {
				continue
			}

			fileName := strings.TrimSuffix(file.Name(), ".tmpl")

			// fileName = strings.TrimPrefix(fileName, "gen/protovalidate/")

			fileName = strings.ReplaceAll(fileName, "___", "/")

			fmt.Fprintf(os.Stderr, "file: %+v\n", fileName)

			var buf bytes.Buffer
			err = file.Execute(&buf, map[string]any{"GoPackageOption": gpkg})
			if err != nil {
				return err
			}

			responseWriter.AddFile(
				filepath.Join(rootPath, "protovalidate", fileName),
				buf.String(),
			)
		}

	}

	return nil
}
