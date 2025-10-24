package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	// File Open
	filename := "words.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Panic("error Open: ", err)
	}

	defer file.Close()

	// Dir Create
	dirname := "chunk_" + strings.TrimSuffix(filename, ".txt")
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.Mkdir(dirname, 0755); err != nil {
			log.Panic("error Create:", err)
		}
	}

	// Reader
	reader := bufio.NewReader(file)
	buffer := make([]byte, 2048)

	// Read Chunk & Write Chunk
	for i := 0; ; i++ {
		n, err := reader.Read(buffer)
		if n > 0 { 
			// Read
			content := string(buffer[:n])
			
			// Write
			chunkFilename := dirname + "/chunk_" + strings.TrimSuffix(filename, ".txt") + "_" + strconv.Itoa(i + 1) + ".txt"
			chunkFile, err := os.Create(chunkFilename)
			if err != nil {
				log.Panic("error Create chunk: ", err)
			}
			if _, err := chunkFile.WriteString(content); err != nil {
				log.Panic("error Writing chunk: ", err)
			}
			chunkFile.Close()
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic("error Read: ", err)
		}
	}
}
