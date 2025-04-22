package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/walteh/protoc-gen-protovalidate/pkg/download"
)

func main() {

	ctx := context.Background()

	languages := []string{"go", "python", "cc", "java"}

	err := os.MkdirAll("gen", 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, language := range languages {
		fmt.Printf("Downloading %s\n", language)
		jsn, tarz, err := download.Download(ctx, language, "latest")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filepath.Join("gen", fmt.Sprintf("protovalidate-%s-latest.json", language)), jsn, 0644)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filepath.Join("gen", fmt.Sprintf("protovalidate-%s-latest.tar.gz", language)), tarz, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
