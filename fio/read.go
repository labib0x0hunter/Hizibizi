package main

import "io"

func Fread(r io.Reader, v ...interface{}) (n int, err error) {

}

func Read(v ...interface{}) (n int, err error) {
	return Fread(In, v...)
}

func main() {

}
