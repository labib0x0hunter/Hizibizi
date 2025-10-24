package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func walkDir(root string, storage map[string][]string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if !d.IsDir() {
			switch ext := filepath.Ext(d.Name()); ext {
			case ".txt":
				storage["TXT_FILE"] = append(storage["TXT_FILE"], d.Name())
			case ".go":
				storage["GO_FILE"] = append(storage["GO_FILE"], d.Name())
			}
		}
		
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

	storage := make(map[string][]string)
	walkDir(rootFolder, storage)

	for fileExtension, files := range storage {
		fmt.Println(fileExtension)
		for _, filename := range files {
			fmt.Println(filename)
		}
		fmt.Println()
	}
}
