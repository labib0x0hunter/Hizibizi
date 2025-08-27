package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strconv"
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
	var val int64 = v.Int()
	w.buf.WriteString(strconv.Itoa(int(val)))
}

func (w *worker) formatString(v reflect.Value) {
	w.buf.WriteString(v.String())
}

func (w *worker) formatStruct(v reflect.Value) {
	w.buf.WriteString("{ ")
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Int:
			w.formatInt(f)
		case reflect.String:
			w.formatString(f)
		case reflect.Struct:
			w.formatStruct(f)
		default:
			return
		}
		if i+1 < v.NumField() {
			w.buf.WriteString(", ")
		}
	}
	w.buf.WriteString(" }")
}

func (w *worker) processFormat(v ...interface{}) {
	for idx, arg := range v {
		vr := reflect.ValueOf(arg)
		switch vr.Kind() {
		case reflect.Int:
			w.formatInt(vr)
		case reflect.String:
			w.formatString(vr)
		case reflect.Struct:
			w.formatStruct(vr)
		default:
			w.free()
			w.buf.WriteString("unknown data type for")
			return
		}
		if idx+1 < len(v) {
			w.buf.WriteByte(' ')
		}
	}
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

type A struct {
	b string
	z int
}

func main() {
	var i int = 10
	j := A{b: "ABC", z: 100}
	Write("Hello %s", "labib", i, j)
	Write("\n")
}
