package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type reader struct {
	reader io.ByteScanner
	buf    buffer
}

var readerPool = sync.Pool{
	New: func() any {
		return new(reader)
	},
}

func newReader(r io.Reader) *reader {
	rdr := readerPool.Get().(*reader)
	if rr, ok := r.(io.ByteScanner); ok {
		rdr.reader = rr
	} else {
		rdr.reader = &byteReader{r: r, unread: false, firstCall: true}
	}

	return rdr
}

func (r *reader) handleError(err string) {
	panic(errors.New(err))
}

func (r *reader) skipDelimater() {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if b != ' ' && b != '\n' {
			r.reader.UnreadByte()
			break
		}
	}
}

func (r *reader) notEof() {
	_, err := r.reader.ReadByte()
	if err != nil || err == io.EOF {
		panic(err)
	}
	r.reader.UnreadByte()
}

func (r *reader) isDigit(b byte) bool {
	for _, x := range "0123456789" {
		if x == rune(b) {
			return true
		}
	}
	return false
}

func (r *reader) isChar(b byte) bool {
	for _, x := range "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if x == rune(b) {
			return true
		}
	}
	return false
}

func (r *reader) readDigit() {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if !r.isDigit(b) {
			r.reader.UnreadByte()
			break
		}
		r.buf.appendByte(b)
	}
}

func (r *reader) readInt() int64 {
	r.skipDelimater()
	r.notEof()

	r.readDigit()
	n, err := strconv.ParseInt(string(r.buf), 10, 64)
	if err != nil {
		// handle error
		panic(err)
	}
	return n
}

func (r *reader) readString() (s string) {
	r.skipDelimater()
	r.notEof()

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if !r.isChar(b) {
			r.reader.UnreadByte()
			break
		}
		r.buf.appendByte(b)
	}
	return string(r.buf)
}

func (r *reader) singleRead(arg reflect.Value) {
	r.buf = r.buf[:0] // clear previous buffer
	arg = arg.Elem()

	switch arg.Kind() {
	case reflect.Int:
		arg.SetInt(int64(r.readInt()))
	case reflect.String:
		arg.SetString(r.readString())
	default:
		r.handleError("unknown data type " + arg.Kind().String())
	}

}

func (r *reader) processRead(v ...interface{}) (n int, err error) {
	for _, arg := range v {
		rv := reflect.ValueOf(arg)
		if rv.Kind() != reflect.Ptr {
			r.handleError("need pointer for " + rv.Type().Name())
			break
		}
		r.singleRead(rv)
		n++
	}
	// Error - handling missing here ...

	// // READ UNTIL EOF or NL
	// for {
	// 	b, err := r.reader.ReadByte()
	// 	if err != nil || err == io.EOF {
	// 		break
	// 	}

	// 	if b == '\n' {
	// 		break
	// 	}
	// }
	return
}

func (r *reader) free() {
	r.buf.truncate()
	r.reader = nil
	readerPool.Put(r)
}

func Fread(r io.Reader, v ...interface{}) (n int, err error) {
	rdr := newReader(r)
	n, err = rdr.processRead(v...)
	rdr.free()
	return
}

func Read(v ...interface{}) (n int, err error) {
	return Fread(os.Stdin, v...)
}

type myInt int

func main() {

	var x int
	var s string

	Read(&x, &s)

	fmt.Println("x: ", x)
	fmt.Println("s: ", s)

	var y myInt
	Read(&y)

	fmt.Println("y: ", y)

}
