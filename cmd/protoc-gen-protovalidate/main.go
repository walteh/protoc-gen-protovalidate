package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

//go:generate go run ./generator

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

	var desc *descriptorpb.FileOptions
	var fdesc protoreflect.FileDescriptor

	for _, fileDescriptor := range fileDescriptors {

		if !strings.HasSuffix(fileDescriptor.Path(), params.GetBufValidateFile()) {
			continue
		}

		descd, ok := fileDescriptor.Options().(*descriptorpb.FileOptions)
		if !ok {
			return fmt.Errorf("file descriptor is not a FileDescriptorProto")
		}
		fdesc = fileDescriptor
		desc = descd
		break
	}

	if desc == nil {
		return nil
	}

	files, err := downloadRemoteFiles(ctx, params.GetLanguage(), params.GetProtoValidateRef())
	if err != nil {
		return err
	}

	switch params.GetLanguage() {
	case "go":
		files, err = GenerateGo(ctx, files, desc, fdesc)
		if err != nil {
			return err
		}
	// case "python":
	// 	files, err = GeneratePython(ctx, files, desc)
	// 	if err != nil {
	// 		return err
	// 	}
	// case "cc":
	// 	files, err = GenerateCC(ctx, files, desc)
	// 	if err != nil {
	// 		return err
	// 	}
	// case "java":
	// 	files, err = GenerateJava(ctx, files, desc)
	// 	if err != nil {
	// 		return err
	// 	}
	default:
		return fmt.Errorf("language not supported: %s", params.GetLanguage())
	}

	for fileName, fileContent := range files {
		responseWriter.AddFile(fileName, fileContent)
	}

	return nil
}

//go:embed gen/protovalidate-*-latest.tar.gz
var localFiles embed.FS

func downloadRemoteFiles(ctx context.Context, language string, ref string) (map[string]string, error) {

	if ref == "_local" {
		localFiles, err := localFiles.ReadFile(filepath.Join("gen", fmt.Sprintf("protovalidate-%s-latest.tar.gz", language)))
		if err != nil {
			return nil, err
		}

		return Untar(bytes.NewReader(localFiles))
	}

	url := fmt.Sprintf("https://github.com/bufbuild/protovalidate-%s/archive/refs/tags/%s.tar.gz", language, ref)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return Untar(resp.Body)
}
