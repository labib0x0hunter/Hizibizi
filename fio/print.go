package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
)

var In io.Reader = os.Stdin
var Out io.Writer = os.Stdout

type buffer struct {
	bytes.Buffer
}

func (b *buffer) appendString(s string) {
	b.WriteString(s)
}

func (b *buffer) appendByte(s byte) {
	b.WriteByte(s)
}

type worker struct {
	buf buffer
	err error
}

func newWorker() *worker {
	return &worker{
		err: nil,
	}
}

func (w *worker) free() {
	w.buf.Reset()
}

func (w *worker) formatInt(v reflect.Value) {
	if v.Kind() != reflect.Int {
		w.err = fmt.Errorf("int=%v", v)
		return
	}
}

func (w *worker) formatString(v reflect.Value) {
	if v.Kind() != reflect.String {
		w.err = fmt.Errorf("string=%v", v)
		return
	}
}

func (w *worker) formatBinary(v reflect.Value) {
	if v.Kind() != reflect.Int {
		w.err = fmt.Errorf("int=%v", v)
		return
	}
}

func (w *worker) formatStruct(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		w.err = fmt.Errorf("struct=%v", v)
		return
	}
}

func (w *worker) processFormat(v ...interface{}) {
}

func Fwrite(w io.Writer, v ...interface{}) (n int, err error) {
	wkr := newWorker()
	wkr.processFormat(v...)
	n, err = w.Write(wkr.buf.Bytes())
	wkr.free()
	return
}

func Write(v ...interface{}) (int, error) {
	return Fwrite(Out, v...)
}

func main() {
	Write("Hello %s", "labib")
	fmt.Println()
	fmt.Printf("%d", "ge")
}
