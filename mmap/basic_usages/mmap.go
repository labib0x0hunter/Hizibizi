package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

// We will store only integer.

var encoder = binary.BigEndian // byte order for encoding

type Index struct {
	file      *os.File    // file metadata and descriptor
	mmap      gommap.MMap // memory-mapped file
	len       uint64      // current filesize
	entrySize uint64      // we will store uint64 only, which is 8byte
	maxSize   uint64
}

func NewIndex(filename string, maxSize int64) *Index {
	// Open file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil
	}

	// Current size of file
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil
	}

	// Truncate to max size
	if err := file.Truncate(maxSize); err != nil {
		file.Close()
		return nil
	}

	// mmap the file
	mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil
	}

	return &Index{
		file:      file,
		mmap:      mmap,
		len:       uint64(info.Size()),
		entrySize: 8,
		maxSize:   uint64(maxSize),
	}
}

// read the number which is stored in offset
func (i *Index) ReadAt(offset uint64) (uint64, error) {
	if offset+i.entrySize > i.maxSize {
		return 0, fmt.Errorf("Overflow")
	}
	num := encoder.Uint64(i.mmap[offset : offset+i.entrySize])
	return num, nil
}

// write an number at offset
func (i *Index) WriteAt(offset, num uint64) error {
	if offset+i.entrySize > i.maxSize {
		return fmt.Errorf("Overflow")
	}
	encoder.PutUint64(i.mmap[offset:offset+i.entrySize], num)

	// For only appending, we will update the len
	// For overriding, the len remains same
	if offset >= i.len {
		i.len += i.entrySize
	}
	return nil
}

// close the mapped file
func (i *Index) Close() error {
	// ensure data leaves mmap and goes to file cache
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	// ensure data reaches disk
	if err := i.file.Sync(); err != nil {
		return err
	}
	// resize file to current size
	if err := i.file.Truncate(int64(i.len)); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	// release the memory
	if err := i.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	// close the file
	if err := i.file.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	filename := "integer.index"
	maxFileSize := 10 * 1024
	store := NewIndex(filename, int64(maxFileSize)) // 10kb
	if store == nil {
		return
	}
	defer store.Close()

	store.WriteAt(0, 10)  // [0:8]
	store.WriteAt(8, 20)  // [8:16]
	store.WriteAt(16, 34) // [16:24]

	//
	num, _ := store.ReadAt(0)
	fmt.Println(num, "stored at offset 0")

	num, _ = store.ReadAt(16)
	fmt.Println(num, "stored at offset 16")

	os.Remove(filename)
}
