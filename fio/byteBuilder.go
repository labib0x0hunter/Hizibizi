package fio

import "unsafe"

// Write at len(buf)
type ByteBuilder struct {
	buf []byte
}

// reset, resize to len 0, but no allocation
func (b *ByteBuilder) Reset() {
	b.buf = b.buf[:0]
}

// bytes
func (b *ByteBuilder) Bytes() []byte {
	return b.buf[:len(b.buf)]
}

// string
// b.buf can change the string..
func (b *ByteBuilder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

func (b *ByteBuilder) Len() int {
	return len(b.buf)
}

// if capacity is less than l + n, then grow by 2 * l + n
// allocate a new slice of capacity
// copy the old buffer to new
// assign new to b.buf
// reslice b.buf by l + n
func (b *ByteBuilder) grow(n int) int {
	l := len(b.buf)
	c := l + n
	if c > cap(b.buf) {
		c *= 2
	}
	temp := make([]byte, c)
	copy(temp, b.buf[:])
	b.buf = temp
	b.buf = b.buf[:l+n]
	return l
}

// if there is enough room for n
// just resize by n + l
func (b *ByteBuilder) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); l+n <= cap(b.buf) {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// append s to b.buf
func (b *ByteBuilder) WriteString(s string) int {
	idx, ok := b.tryGrowByReslice(len(s))
	if !ok {
		idx = b.grow(len(s))
	}
	return copy(b.buf[idx:], s)
}

// append s to b.buf
func (b *ByteBuilder) WriteBytes(s []byte) int {
	idx, ok := b.tryGrowByReslice(len(s))
	if !ok {
		idx = b.grow(len(s))
	}
	return copy(b.buf[idx:], s)
}

// append s to b.buf
func (b *ByteBuilder) WriteByte(s byte) int {
	idx, ok := b.tryGrowByReslice(1)
	if !ok {
		idx = b.grow(1)
	}
	b.buf[idx] = s
	return 1
}

// func main() {

// 	var s ByteBuilder

// 	s.WriteString("hello")
// 	// println(len(s.buf), cap(s.buf))
// 	// fmt.Println(s.String())

// 	s.WriteString("bb")
// 	// println(len(s.buf), cap(s.buf))
// 	// fmt.Println(s.String())

// 	s.WriteString("abcdefg")
// 	// println(len(s.buf), cap(s.buf))
// 	// fmt.Println(s.String())

// 	s.WriteByte('k')
// 	// println(len(s.buf), cap(s.buf))
// 	// fmt.Println(s.String())
// }
