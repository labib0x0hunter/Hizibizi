package main

import (
	"io"
	"reflect"
)

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

type worker struct {
	buf buffer
}

func newWorker() *worker {
	return &worker{}
}

func (w *worker) free() {
	w.buf = buffer{}
}

func (w *worker) formatInt(v reflect.Value) {
	var val int64 = v.Int()
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

func (w *worker) formatString(v reflect.Value) {
	w.buf.appendString(v.String())
}

func (w *worker) formatStruct(v reflect.Value) {
	w.buf.appendString("{ ")
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
			w.buf.appendString(", ")
		}
	}
	w.buf.appendString(" }")
}

func (w *worker) processWrite(v ...interface{}) {
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
			w.buf.appendString("unknown data type for")
			return
		}
		if idx+1 < len(v) {
			w.buf.appendByte(' ')
		}
	}
}

func Fwrite(w io.Writer, v ...interface{}) (n int, err error) {
	wkr := newWorker()
	wkr.processWrite(v...)
	n, err = w.Write(wkr.buf)
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
