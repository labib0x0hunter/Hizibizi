package fio

import (
	"io"
	"os"
	"unsafe"
)

const bufSize = 3

type Writer struct {
	buf []byte
	w   io.Writer // os.Stdout
	n   int       // how much stored
}

func BufferWriter() *Writer {
	return &Writer{
		buf: make([]byte, bufSize),
		w:   os.Stdout,
		n:   0,
	}
}

func (b *Writer) IsWritten() bool {
	return b.n == 0
}

func (b *Writer) Availabel() int {
	return bufSize - b.n
}

func (b *Writer) Flush() error {
	if b.IsWritten() {
		return nil
	}

	n, err := b.w.Write(b.buf[:b.n])
	if err != nil {
		return err
	}

	b.n -= n
	copy(b.buf, b.buf[n:])
	return nil
}

func (b *Writer) WriteBytes(s []byte) (n int, err error) {
	for b.Availabel() < len(s) && err == nil {
		if b.IsWritten() {
			n, err = b.w.Write(s[:bufSize])
			s = s[n:]
		} else {
			l := copy(b.buf[b.n:], s)
			s = s[l:]
			b.n += l
			err = b.Flush()
		}
	}

	copy(b.buf[b.n:], s)
	b.n += len(s)
	return
}

func (b *Writer) WriteString(s string) (n int, err error) {
	if s == "" {
		return 0, nil
	}
	return b.WriteBytes(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// func main() {
// 	w := BufferWriter()
// 	w.WriteString("hello")
// 	w.WriteString("\n")

// 	w.Flush()
// }
