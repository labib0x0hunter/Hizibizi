package main

import (
	"errors"
	"fmt"
	"time"
)

type Time struct {}

func (t Time) Name() string {
	return "Running time"
}

func (t Time) Run(args []string) error {
	if len(args) != 0 {
		return errors.New("too many argument")
	}
	fmt.Println(time.Now().Format("2006-01-02"))
	return nil
}