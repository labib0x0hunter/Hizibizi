package main

import (
	"fmt"
	"strconv"
)

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func mult(a, b int) int {
	return a * b
}

func div(a, b int) int {
	return a / b
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	var onMap = map[string]func(int, int) int{
		"+": add,
		"-": sub,
		"*": mult,
		"/": div,
	}

	expressions := [][]string{
		{"2", "+", "3"},
		{"2", "-", "3"},
		{"2", "*", "3"},
		{"2", "/", "3"},
		{"2", "%", "3"},
	}

	for _, exp := range expressions {
		a, err := strconv.Atoi(exp[0])
		Panic(err)
		b, err := strconv.Atoi(exp[2])
		Panic(err)
		
		op := exp[1]
		if funcOp, ok := onMap[op]; ok {
			fmt.Println(funcOp(a, b))
			continue
		}

		fmt.Println("unknown operation")
	}

}
