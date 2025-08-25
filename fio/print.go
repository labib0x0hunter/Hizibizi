package main

import (
	"bytes"
	"io"
	"os"
)

var In io.Reader = os.Stdin
var Out io.Writer = os.Stdout

type buffer struct {
	bytes.Buffer
}

func (b *buffer) append(s string) {
	b.WriteString(s)
}

type worker struct {
	buf buffer
}

func (w *worker) free() {
	w.buf.Reset()
}

func newWorker() *worker {
	return &worker{}
}

func Fwrite(w io.Writer, format string, v ...interface{}) (n int, err error) {
	wkr := newWorker()
	for _, val := range v {
		switch x := val.(type) {
		case string:
			wkr.buf.append(x)
		}
	}
	n, err = w.Write(wkr.buf.Bytes())
	wkr.free()
	return
}

func Write(format string, v ...interface{}) (int, error) {
	return Fwrite(Out, format, v...)
}

func main() {
	Write("", "HELLO WORLD\n")
}
