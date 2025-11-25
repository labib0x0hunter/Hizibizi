package fio

import (
	"fmt"
	"io"
	"os"
)

// const bufSize = 100

type Reader struct {
	buf []byte
	r   io.Reader
	n   int // how much read
	// rp, wp int 	// read write pointer
}

func BufferReader() *Reader {
	return &Reader{
		buf: make([]byte, bufSize),
		r:   os.Stdin,
	}
}

func (r *Reader) noRead() bool {
	return r.n == 0
}

func (r *Reader) Avalable() int {
	return bufSize - r.n
}

func (r *Reader) read(limit int) {
	n, _ := r.r.Read(r.buf[limit:])
	r.n += n
}

func (r *Reader) ReadBytes(p []byte) (n int, err error) {
	if r.noRead() {
		r.read(0)
	}
	n = copy(p, r.buf[:])
	r.buf = r.buf[n:]
	r.n -= n
	return
}

func (r *Reader) Read(v ...interface{}) {

}

func (r *Reader) isDigit(b byte) bool {
	for _, d := range "0123456789" {
		if d == rune(b) {
			return true
		}
	}
	return false
}

func (r *Reader) ReadInt(n *int) (err error) {
	if r.noRead() {
		r.read(0)
	}

	lastDigit := 0
	for i := 0; i < cap(r.buf); i++ {
		if i < len(r.buf) {
			if r.isDigit(r.buf[i]) {
				lastDigit++
				continue
			} else {
				break
			}
		} else {
			r.read(len(r.buf))
			i--
		}
	}

	for i := 0; i < lastDigit; i++ {
		*n = (*n * 10) + (int(r.buf[i] - '0'))
	}

	return
}

func main() {
	r := BufferReader()

	// b := make([]byte, 1)

	// r.ReadBytes(b)

	// fmt.Println(b)

	var n int
	r.ReadInt(&n)

	fmt.Println(n)
}
