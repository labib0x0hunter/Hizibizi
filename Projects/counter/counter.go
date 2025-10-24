package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func writeFile(filename string, cnt int) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("error opening file : ", err)
	}
	defer file.Close()
	if _, err := file.WriteString(strconv.Itoa(cnt)); err != nil {
		log.Println("error writing file : ", err)
	}
}

func readFile(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("error opening file : ", err)
		return 0, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("error reading file : ", err)
		return 0, err
	}
	num, err := strconv.Atoi(string(data))
	if err != nil {
		log.Println("error converting string :", err, string(data))
		return 0, err
	}
	return num, nil
}

func main() {

	filename := "counter.db"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		file.WriteString("0")
		file.Close()
	}

	if len(os.Args) < 2 {
		fmt.Println("go run counter.go [inc|get|reset]")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "inc":
		count, err := readFile(filename)
		if err == nil {
			writeFile(filename, count + 1)
		}
	case "get":
		count, err := readFile(filename)
		if err == nil {
			fmt.Println(count)
		}
	case "reset":
		writeFile(filename, 0)
	default:
		fmt.Println("Unknown Command")
	}
}
