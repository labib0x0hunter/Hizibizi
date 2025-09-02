package main

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

func (b *byteReader) ReadByte() (byte, error) {
	var bb byte
	var err error
	b.firstCall = false

	if b.unread {
		bb = b.prevBuf
		b.unread = false
		return bb, err
	}

	n, err := io.ReadFull(b.r, b.buf[:]) // read 1 byte
	if n != 1 {
		return bb, errors.New("longer")
	}
	bb = b.buf[0]
	b.prevBuf = bb
	return bb, err
}

func (b *byteReader) UnreadByte() error {
	var err error
	if b.firstCall || b.unread {
		return errors.New("no byte to unread()")
	}
	b.unread = true
	return err
}
