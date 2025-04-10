package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func GenerateGo(ctx context.Context, files map[string]string, opts *descriptorpb.FileOptions, desc protoreflect.FileDescriptor) (map[string]string, error) {
	// take allthe .go files in the main dir, and in the resolve and cid dirs and process them

	var replacements = map[string]string{
		"github.com/bufbuild/protovalidate-go":                                    opts.GetGoPackage() + "/protovalidate",
		"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate": opts.GetGoPackage(),
	}

	filed := map[string]string{}

	rootPath := strings.TrimSuffix(desc.Path(), "validate.proto")

	for filePath, fileContent := range files {

		if strings.HasSuffix(filePath, "_test.go") {
			continue
		}

		// only grab files matching "*.go|resolve/*.go|cid/*.go"

		glob := "{*,resolve/*,cel/*}.go"

		// print the file apth

		fmt.Fprintf(os.Stderr, "filePath: %s\n", filePath)

		matches, err := doublestar.PathMatch(glob, filePath)
		if err != nil {
			return nil, err
		}

		if !matches {
			continue
		}

		content := string(fileContent)

		for old, new := range replacements {
			content = strings.ReplaceAll(content, old, new)
		}

		filed[filepath.Join(rootPath, "protovalidate", filePath)] = content
	}

	if len(filed) == 0 {
		return nil, errors.New("no files to process")
	}

	return filed, nil
}
