package fio

import (
	"io"
	"os"
	"reflect"
	"sync"
)

type writer struct {
	buf buffer
}

var writerPool = sync.Pool{
	New: func() any {
		return new(writer)
	},
}

// newWriter retrieves a *writer instance from the workerPool.
// It is used to obtain a reusable writer object for performing write operations.
func newWriter() *writer {
	w := writerPool.Get().(*writer)
	return w
}

// free releases resources held by the writer instance.
// It truncates the internal buffer and returns the writer to the worker pool for reuse.
func (w *writer) free() {
	w.buf.truncate()
	writerPool.Put(w)
}

// formatInt converts an integer value to its decimal string representation
// and appends the result to the writer's buffer. It handles zero explicitly
// by appending '0'. Negative values are not handled.
func (w *writer) formatInt(val int) {
	// var val int64 = v.Int()
	if val == 0 {
		w.buf.appendByte('0')
		return
	}
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

// formatString appends the provided string v to the writer's internal buffer.
func (w *writer) formatString(v string) {
	w.buf.appendString(v)
}

// formatStruct formats the fields of a struct value and appends the result to the writer's buffer.
// It handles fields of type int, string, and nested structs recursively.
// For unsupported field types, it appends a placeholder indicating an unknown type and returns early.
// Fields are separated by commas and enclosed in curly braces.
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


// processWrite formats and writes each argument in the provided slice v to the writer's buffer.
// It handles arguments of type int, string, and struct, formatting them accordingly.
// For unsupported types, it appends a placeholder indicating an unknown type and returns early.
// Arguments are separated by a space in the buffer.
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

// Fwrite writes the provided values to the given io.Writer.
// It processes the input values, serializes them into a buffer, and writes the buffer to w.
// Returns the number of bytes written and any error encountered during the write operation.
func Fwrite(w io.Writer, v ...interface{}) (n int, err error) {
	wkr := newWriter()
	wkr.processWrite(v...)
	n, err = w.Write(wkr.buf)
	wkr.free()
	return
}

// Write writes the provided values to the standard output using Fwrite.
// It returns the number of bytes written and any error encountered.
// The function accepts a variadic number of interface{} arguments to support
// writing multiple values at once.
func Write(v ...interface{}) (int, error) {
	return Fwrite(os.Stdout, v...)
}

// type A struct {
// 	b string
// 	z int
// }

// type B struct {
// 	a A
// 	b string
// }

// func main() {
// 	var i int = 102222
// 	j := A{b: "ABC", z: 100}
// 	k := B{a: j, b: "HELLO"}
// 	Write("Hello %s", "labib", i, j, k)
// 	Write("\n")
// }
