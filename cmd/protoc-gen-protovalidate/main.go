package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protoplugin"
	"github.com/walteh/protoc-gen-protovalidate/pkg/download"
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

	jout, tout, err := downloadRemoteFiles(ctx, params.GetLanguage(), params.GetProtoValidateRef())
	if err != nil {
		return err
	}

	files := make(map[string]string)

	switch params.GetLanguage() {
	case "go":
		files, err = GenerateGo(ctx, tout, desc, fdesc, jout)
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

//go:embed gen/protovalidate-*-latest*
var localFiles embed.FS

func downloadRemoteFiles(ctx context.Context, language string, ref string) (map[string]any, map[string]string, error) {
	var jbytes []byte
	var tbytes []byte
	var err error
	if ref == "_local" {

		jbytes, err = localFiles.ReadFile(filepath.Join("gen", fmt.Sprintf("protovalidate-%s-latest.json", language)))
		if err != nil {
			return nil, nil, err
		}

		tbytes, err = localFiles.ReadFile(filepath.Join("gen", fmt.Sprintf("protovalidate-%s-latest.tar.gz", language)))
		if err != nil {
			return nil, nil, err
		}

	} else {
		jbytes, tbytes, err = download.Download(ctx, language, ref)
		if err != nil {
			return nil, nil, err
		}

	}

	var jout map[string]interface{}
	err = json.Unmarshal(jbytes, &jout)
	if err != nil {
		return nil, nil, err
	}

	tout, err := Untar(bytes.NewReader(tbytes))
	if err != nil {
		return nil, nil, err
	}

	return jout, tout, nil

}
