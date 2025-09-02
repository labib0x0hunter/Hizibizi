package main

import (
	"io"
	"os"
	"reflect"
	"sync"
)

type writer struct {
	buf buffer
}

var workerPool = sync.Pool{
	New: func() any {
		return new(writer)
	},
}

func newWriter() *writer {
	w := workerPool.Get().(*writer)
	return w
}

func (w *writer) free() {
	w.buf.truncate()
	workerPool.Put(w)
}

// formats int, digit by digit
// append to buffer
func (w *writer) formatInt(val int) {
	// var val int64 = v.Int()
	var inbuf [20]byte
	i := len(inbuf)

	for val > 0 {
		i--
		nxt := val / 10
		inbuf[i] = byte('0' + val - nxt*10)
		val = nxt
	}
	w.buf.appendBytes(inbuf[i:])
}

// format string and append to buffer
func (w *writer) formatString(v string) {
	w.buf.appendString(v)
}

// format struct, and append to buffer
func (w *writer) formatStruct(v reflect.Value) {
	w.buf.appendString("{ ")
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Int:
			w.formatInt(int(f.Int()))
		case reflect.String:
			w.formatString(f.String())
		case reflect.Struct:
			w.formatStruct(f)
		default:
			w.buf.appendString(" %!unknown% ")
			return
		}
		if i+1 < v.NumField() {
			w.buf.appendString(", ")
		}
	}
	w.buf.appendString(" }")
}

// detect type by reflection, then format them
// other than string, int and struct -> error
func (w *writer) processWrite(v ...interface{}) {
	for idx, arg := range v {
		vr := reflect.ValueOf(arg)
		switch vr.Kind() {
		case reflect.Int:
			w.formatInt(int(vr.Int()))
		case reflect.String:
			w.formatString(vr.String())
		case reflect.Struct:
			w.formatStruct(vr)
		default:
			w.buf.appendString(" %!unknown% ")
			return
		}
		if idx+1 < len(v) {
			w.buf.appendByte(' ')
		}
	}
}

// Accepts an io.Writer and any arguments
func Fwrite(w io.Writer, v ...interface{}) (n int, err error) {
	wkr := newWriter()
	wkr.processWrite(v...)
	n, err = w.Write(wkr.buf)
	wkr.free()
	return
}

func Write(v ...interface{}) (int, error) {
	return Fwrite(os.Stdout, v...)
}

type A struct {
	b string
	z int
}

type B struct {
	a A
	b string
}

func main() {
	var i int = 102222
	j := A{b: "ABC", z: 100}
	k := B{a: j, b: "HELLO"}
	Write("Hello %s", "labib", i, j, k)
	Write("\n")
}
