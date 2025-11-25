package fio

import (
	"io"
	"os"
)

var (
	In  io.Reader = os.Stdin	// the standard input reader.
	Out io.Writer = os.Stdout	// the standard output writer.
)

// buffer is a type alias for a byte slice, used for storing buffered data.
type buffer []byte

// appendString appends the contents of string s to the buffer.
func (b *buffer) appendString(s string) {
	*b = append(*b, s...)
}

// appendBytes appends the contents of byte slice s to the buffer.
func (b *buffer) appendBytes(s []byte) {
	*b = append(*b, s...)
}

// appendByte appends a single byte s to the buffer.
func (b *buffer) appendByte(s byte) {
	*b = append(*b, s)
}

// truncate clears the buffer by setting its length to zero;
// the underlying array is not reallocated and its capacity remains unchanged.
func (b *buffer) truncate() {
	*b = (*b)[:0]
}