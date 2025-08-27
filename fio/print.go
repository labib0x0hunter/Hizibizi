package main

import (
	"bytes"
	"fmt"
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
	if v.Kind() != reflect.Struct {
		w.err = fmt.Errorf("struct=%v", v)
		return
	}
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
		if idx + 1 < len(v) {
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

func main() {
	var i int = 10
	Write("Hello %s", "labib", i)
	Write("\n")
}
