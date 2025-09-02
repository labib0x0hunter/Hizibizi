package main

import (
	"io"
	"os"
)

// std input/output
var In io.Reader = os.Stdin
var Out io.Writer = os.Stdout

// to store buffer
type buffer []byte

func (b *buffer) appendString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) appendBytes(s []byte) {
	*b = append(*b, s...)
}

func (b *buffer) appendByte(s byte) {
	*b = append(*b, s)
}

func (b *buffer) truncate() {
	*b = (*b)[0:]
}
