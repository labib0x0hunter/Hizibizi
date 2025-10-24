package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func walkDir(path string, depth int) {
	if depth == 0 {
		fmt.Println(path)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, filename := range files {
		fmt.Println(strings.Repeat("   ", depth), filename.Name())
		if filename.IsDir() {
			walkDir(filepath.Join(path, filename.Name()), depth + 1)
		}
	}
}

func walkDirFilepathPackage(root string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			panic(err)
		}
		depth := strings.Count(relPath, string(os.PathSeparator))
		fmt.Println(strings.Repeat("   ", depth), d.Name())
		return nil
	})
}

func main() {

	currentFolder, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootFolder := filepath.Join(currentFolder, "../../")
	rootFolder = filepath.Clean(rootFolder)

	// walkDir(rootFolder, 0)
	walkDirFilepathPackage(rootFolder)

}