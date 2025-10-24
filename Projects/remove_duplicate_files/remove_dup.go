package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getFiles(path string) []string {
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, filename := range dir {
		if strings.HasSuffix(filename.Name(), ".txt") {
			files = append(files, filename.Name())
		}
	}
	return files
}

func getContent(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return content
}

func generateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

func getHash(path string) string {
	content := getContent(path)
	return generateHash(content)
}

func main() {
	dirName := "tests/"
	files := getFiles(dirName)
	seen := make(map[string]struct{})
	for _, filename := range files {
		path := filepath.Join(dirName, filename)
		hash := getHash(path)
		// hash := efficientHashGenerator(path)
		if _, ok := seen[hash]; ok {
			if err := os.Remove(path); err != nil {
				panic(err)
			}
			continue
		}
		seen[hash] = struct{}{}
	}
}