package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dirs = []string{
	"../..",
	"../../cel",
	"../../resolve",
}

var replacements = map[string]string{
	"github.com/bufbuild/protovalidate-go":                                    "[[[[[[[[.GoPackageOption]]]]]]]]/protovalidate",
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate": "[[[[[[[[.GoPackageOption]]]]]]]]",
}

func main() {

	// grab all the .go files in the ../../* ../../cel/* and ../../resolve/* directories
	// write them to ./embed/protovalidate/*

	os.RemoveAll("./templates")

	os.MkdirAll("./templates", 0755)

	files := []string{}

	for _, dir := range dirs {
		g, err := filepath.Glob(fmt.Sprintf("%s/*.go", dir))
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, g...)
	}

	for _, file := range files {
		// do not include test files
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		outFile := strings.Replace(file, "../../", "", 1) + ".tmpl"

		dataStr := string(data)

		// escape the data as a go template
		for k, v := range replacements {
			dataStr = strings.ReplaceAll(dataStr, k, v)
		}

		outFile = strings.ReplaceAll(outFile, "/", "___")

		err = os.WriteFile(filepath.Join("templates", outFile), []byte(dataStr), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

}
