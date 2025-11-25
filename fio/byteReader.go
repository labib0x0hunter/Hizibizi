package fio

import (
	"errors"
	"io"
)

type byteReader struct {
	r         io.Reader
	buf       [1]byte
	unread    bool
	prevBuf   byte
	firstCall bool
}

// ReadByte reads and returns a single byte from the underlying reader.
// If the unread flag is set, it returns the previously read byte and resets the flag.
// Otherwise, it reads one byte from the reader, stores it for potential unread operations,
// and returns the byte. If reading fails or does not return exactly one byte, an error is returned.
func (b *byteReader) ReadByte() (byte, error) {
	var bb byte
	var err error
	b.firstCall = false

	if b.unread {
		bb = b.prevBuf
		b.unread = false
		return bb, nil
	}

	n, err := io.ReadFull(b.r, b.buf[:]) // read 1 byte
	if n != 1 {
		return bb, errors.New("failed to read 1 byte")
	}
	bb = b.buf[0]
	b.prevBuf = bb
	return bb, err
}

// UnreadByte marks the last read byte as unread, allowing it to be read again on the next read operation.
// Returns an error if there is no byte to unread, such as when called before any read or if a byte is already unread.
func (b *byteReader) UnreadByte() error {
	if b.firstCall || b.unread {
		return errors.New("no byte to unread()")
	}
	b.unread = true
	return nil
}
