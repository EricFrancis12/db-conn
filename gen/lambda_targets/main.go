package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"db-conn/pkg"
)

func main() {
	var (
		filePath   = filepath.Join("../../", pkg.DefaultFilePath)
		outputPath = "../../lambda/gen.go"
	)

	connStrs, err := pkg.ReadToConnStrs(filePath)
	if err != nil {
		log.Fatal(err)
	}

	joined := ""
	for _, s := range connStrs {
		joined += fmt.Sprintf("`%s`,\n", s)
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
