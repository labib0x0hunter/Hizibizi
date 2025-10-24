package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func deleteFile(path string, dryRun bool) {
	if dryRun {
		fmt.Println("[dryRun]", path, "is to be removed")
	} else {
		os.Remove(path)
		log.Println(path, "removed")
	}
}

func walkDir(path string, dryRun bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, filename := range files {
		fullPath := filepath.Join(path, filename.Name())
		if filename.IsDir() {
			walkDir(fullPath, dryRun)
		} else {
			info, err := filename.Info()
			if err != nil {
				panic(err)
			}
			if info.Size() != 0 {
				continue
			}

			// Empty files
			deleteFile(fullPath, dryRun)
		}
	}

	// Empty Folder
	remainFiles, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if len(remainFiles) != 0 {
		return
	}

	deleteFile(path, dryRun)
}

func main() {

	currentFolder, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootFolder := filepath.Join(currentFolder, "../../")
	rootFolder = filepath.Clean(rootFolder)

	dryRun := true // Set to false if you want to delete
	walkDir(rootFolder, dryRun)
}
