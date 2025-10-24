package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 4 {
		panic("Usage: go run calculator.go 10 '+' 20")
	}

	num1, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("")
	}
	num2, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("")
	}
	op := os.Args[2]

	var ans int

	switch op {
	case "+":
		ans = num1 + num2
	case "-":
		ans = num1 - num2
	case "/":
		ans = num1 / num2
	case "*":
		ans = num1 * num2
	}

	fmt.Println(ans)

}