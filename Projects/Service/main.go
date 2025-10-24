package main

import "fmt"

type Logger struct {}

func (l Logger) Log(msg string) {
	fmt.Println("LOG:", msg)
}

type Service struct {
	Logger
	ServiceName string
}

func main() {

	authService := Service{
		ServiceName: "Auth",
	}

	authService.Log("Hello")

}