package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	c := Calc{}
	fmt.Println(GetName(c))
	if err := Runnable(c, []string{"10", "+", "20"}); err != nil {
		log.Println(err)
	}

	d := Hello{}
	fmt.Println(GetName(d))
	if err := Runnable(d, []string{"labib"}); err != nil {
		log.Println(err)
	}

	e := Time{}
	fmt.Println(GetName(e))
	if err := Runnable(e, []string{}); err != nil {
		log.Println(err)
	}

	e = Time{}
	fmt.Println(GetName(e))
	if err := Runnable(e, []string{"Labib", "Al", "Faisal"}); err != nil {
		log.Println(err)
	}

	var commands = map[string]Command{
		"calc":  Calc{},
		"hello": Hello{},
		"time":  Time{},
	}

	if cmd, ok := commands[os.Args[1]]; ok {
		err := Runnable(cmd, os.Args[2:])
		if err != nil {
			log.Println("Error:", err)
		}
	} else {
		fmt.Println("Unknown command:", os.Args[1])
	}

}
