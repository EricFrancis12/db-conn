package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindAncestorDirWith(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found")
		}
		dir = parent
	}
}

func ReadToConnStrs(filePath string) ([]string, error) {
	b, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	return fmtLines(lines), nil
}

func readFile(filePath string) ([]byte, error) {
	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' not found", filePath)
	} else if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	} else if stat.IsDir() {
		return nil, fmt.Errorf("'%s' needs to be a path to a file (found dir)", filePath)
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
	}

	return b, nil
}

func fmtLines(lines []string) []string {
	result := []string{}
	for _, line := range lines {
		s := strings.TrimSpace(line)

		// ignore empty lines & comments
		if s == "" || s[0] == '#' {
			continue
		}

		result = append(result, s)
	}
	return result
}
