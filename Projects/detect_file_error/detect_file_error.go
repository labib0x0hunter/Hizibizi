package main

import (
	"errors"
	"fmt"
	"os"
)

type ReadError struct {
	Filename string
	err      error
}

func (re *ReadError) Error() string {
	return fmt.Sprintf("failed to open %v", re.err)
}

// func (re *ReadError) Unwrap() error {
// 	return re.err
// }

func OpenFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return &ReadError{
			Filename: filename,
			err: err,
		}
	}
	defer file.Close()
	return nil
}

func main() {

	err := OpenFile("a.txt")
	if err == nil {

	}

	// err != nil
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("sentinel", err)
	}
	fmt.Println()

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		fmt.Println("path error")
		fmt.Println("OP: ", pathErr.Op)
	}
	fmt.Println()

	var readErr *ReadError
	if errors.As(err, &readErr) {
		fmt.Println("read error")
		fmt.Println("OP: ", readErr.Filename)
	}
	fmt.Println()
}
