package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"

	"db-conn/pkg"
)

func main() {
	rootPath, err := pkg.FindAncestorDirWith("go.mod")
	if err != nil {
		log.Fatalf("error getting root dir path: %v", err)
	}

	var (
		filePath   = filepath.Join(rootPath, pkg.DefaultFilePath)
		outputPath = filepath.Join(rootPath, "lambda/gen.go")
	)

	connStrs, err := pkg.ReadToConnStrs(filePath)
	if err != nil {
		log.Fatal(err)
	}

	joined := ""
	for _, s := range connStrs {
		safe := strings.ReplaceAll(s, "`", "` + \"`\" + `")
		joined += fmt.Sprintf("`%s`,\n", safe)
	}

	goCode := fmt.Sprintf(
		`package main
		
		var targets = [%d]string{
			%s
		}`,
		len(connStrs),
		joined,
	)

	formattedCode, err := format.Source([]byte(goCode))
	if err != nil {
		log.Fatalf("error formating code: %v", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("error creating file %s: %v", outputPath, err)
	}
	defer out.Close()

	if _, err := out.Write([]byte(formattedCode)); err != nil {
		log.Fatalf("error writing to file %s: %v", outputPath, err)
	}

	fmt.Printf("Successfully generated '%s'\n", outputPath)
}
