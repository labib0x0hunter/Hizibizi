package fio

import (
	"errors"
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

// newReader returns a *reader instance initialized with the provided io.Reader.
// If the input implements io.ByteScanner, it is used directly; otherwise, it wraps
// the reader in a byteReader to provide ByteScanner functionality. The reader
// instance is retrieved from a pool for efficient reuse.
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

// skipDelimater advances the reader past any consecutive space (' ') or newline ('\n') 
// characters. It stops at the first non-delimiter byte, which is then unread so that 
// subsequent reads start from this byte. If the end of the input is reached, the function returns.
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

// notEof checks if the underlying reader has not reached EOF.
// It attempts to read a byte and panics if an error occurs or EOF is encountered.
// If successful, it unreads the byte to restore the reader's position.
func (r *reader) notEof() {
	_, err := r.reader.ReadByte()
	if err != nil || err == io.EOF {
		panic(err)
	}
	r.reader.UnreadByte()
}

// isDigit checks if the provided byte represents an ASCII digit ('0' to '9').
// It returns true if the byte is a digit, and false otherwise.
func (r *reader) isDigit(b byte) bool {
	for _, x := range "0123456789" {
		if x == rune(b) {
			return true
		}
	}
	return false
}

// isChar checks if the given byte represents an alphanumeric character (0-9, a-z, or A-Z).
func (r *reader) isChar(b byte) bool {
	for _, x := range "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if x == rune(b) {
			return true
		}
	}
	return false
}

// readDigit reads consecutive digit bytes from the underlying reader and appends them to the buffer.
// It stops reading when a non-digit byte is encountered or when the end of the input is reached.
// If a non-digit byte is read, it is pushed back to the reader for future processing.
// Panics if an unexpected read error occurs.
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

// readInt reads an integer value from the input buffer, handling delimiters and end-of-file checks.
// It parses the buffered digits as a base-10 int64 and panics if parsing fails.
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

// readString reads a sequence of valid characters from the underlying reader,
// skipping any delimiter at the start. It returns the read string up to the first
// non-character byte or EOF. If an error other than EOF occurs during reading,
// the function panics. The function uses an internal buffer to accumulate bytes
// and returns the resulting string.
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

// singleRead reads a single value from the underlying buffer and sets it to the provided reflect.Value.
// It supports reading values of type int and string. For int, it reads an integer from the buffer and sets it.
// For string, it reads a string from the buffer and sets it. If the value's kind is not supported,
// it calls handleError with an appropriate message.
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

// processRead reads data into the provided pointer arguments using reflection.
// Each argument must be a pointer; otherwise, an error is handled internally.
// The function increments the count for each successful read and returns the total number of reads performed.
// Error handling for the overall read process is currently missing.
// The function is intended to be used for reading multiple values in sequence.
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
	return
}

// free releases resources held by the reader instance.
// It truncates the internal buffer, sets the underlying reader to nil,
// and returns the reader to the pool for reuse.
func (r *reader) free() {
	r.buf.truncate()
	r.reader = nil
	readerPool.Put(r)
}

// Fread reads data from the provided io.Reader into the variables specified by v.
// It returns the number of items successfully read and any error encountered.
// The function uses an internal reader to process the read operation and ensures resources are freed after use.
func Fread(r io.Reader, v ...interface{}) (n int, err error) {
	rdr := newReader(r)
	n, err = rdr.processRead(v...)
	rdr.free()
	return
}

// Read reads input from the standard input (os.Stdin) and stores the values
// into the provided variables v. It returns the number of items successfully
// read and any error encountered during the read operation.
// The function delegates the actual reading to Fread.
func Read(v ...interface{}) (n int, err error) {
	return Fread(os.Stdin, v...)
}

// type myInt int

// func main() {

// 	var x int
// 	var s string

// 	Read(&x, &s)

// 	fmt.Println("x: ", x)
// 	fmt.Println("s: ", s)

// 	var y myInt
// 	Read(&y)

// 	fmt.Println("y: ", y)

// }
