package main

import (
	"errors"
	"fmt"
)

type Hello struct {}

func (h Hello) Name() string {
	return "Running hello"
}

func (h Hello) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("too many argument")
	}
	fmt.Println("Hello", args[0])
	return nil
}