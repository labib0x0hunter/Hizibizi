package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Calc struct {}

func (c Calc) Name() string {
	return "Running calc"
}

func (c Calc) Run(args []string) error {
	if len(args) < 3 {
		return errors.New("not enough argument")
	}

	a, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	b, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	switch args[1] {
	case "+":
		fmt.Println(a + b)
	case "*":
		fmt.Println(a * b)
	default:
		fmt.Println("Unknown cmd")
	}
	return nil
}