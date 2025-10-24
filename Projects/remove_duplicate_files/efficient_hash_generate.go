package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Read file chunk by chunk and generate hash
func efficientHashGenerator(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // 4kb buffer
	hasher := sha256.New()

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			// Hash
			if _, err := hasher.Write(buffer[:n]); err != nil {
				panic(err)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	return hex.EncodeToString(hasher.Sum(nil))
}