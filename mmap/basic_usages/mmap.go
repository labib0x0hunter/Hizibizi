package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

// This program demonstrates how to persist and read fixed-size (8 bytes) uint64 integers
// using memory-mapped files for high-performance I/O and minimal memory overhead.

var encoder = binary.BigEndian // consistent byte order(Big Endian) for binary encoding

type Index struct {
	file      *os.File    // file metadata and descriptor
	mmap      gommap.MMap // memory-mapped view of the file as []byte
	len       uint64      // current filesize
	entrySize uint64      // each entry is uint64 (8 bytes)
	maxSize   uint64      // maximum filesize
}

// NewIndex creates and memory-maps the file with given filename and maxSize.
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

	// Expand file to max size
	if err := file.Truncate(maxSize); err != nil {
		file.Close()
		return nil
	}

	// mmap the file with read/write access, shared among processes
	mmap, err := gommap.Map(
		file.Fd(),                          // File descriptor of the file
		gommap.PROT_READ|gommap.PROT_WRITE, // Memory protection (read & write access)
		gommap.MAP_SHARED,                  // Mapping mode (shared memory between processes), used for persistency
	)
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

// ReadAt reads a uint64 integer from the given byte offset.
func (i *Index) ReadAt(offset uint64) (uint64, error) {
	if offset+i.entrySize > i.maxSize {
		return 0, fmt.Errorf("Overflow")
	}
	num := encoder.Uint64(i.mmap[offset : offset+i.entrySize])
	return num, nil
}

// WriteAt writes a uint64 integer to the given byte offset.
// It updates the logical size if appending new data.
func (i *Index) WriteAt(offset, num uint64) error {
	if offset+i.entrySize > i.maxSize {
		return fmt.Errorf("Overflow")
	}
	encoder.PutUint64(i.mmap[offset:offset+i.entrySize], num)

	// For only appending, we will update the len
	// For overwriting, the len remains same
	if offset >= i.len {
		i.len += i.entrySize
	}
	return nil
}

// Close flushes changes, unmaps memory, truncates unused file space, and closes the file.
func (i *Index) Close() error {
	// sync mmap to file cache
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	// ensure file cache is flushed to disk
	if err := i.file.Sync(); err != nil {
		return err
	}
	// shrink file to current size
	if err := i.file.Truncate(int64(i.len)); err != nil {
		return err
	}
	// ensure file cache is flushed to disk
	if err := i.file.Sync(); err != nil {
		return err
	}
	// unmap memory
	if err := i.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	// close file
	if err := i.file.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	filename := "integer.index"
	maxFileSize := 10 * 1024 // 10kb
	
	store := NewIndex(filename, int64(maxFileSize))
	if store == nil {
		return
	}
	defer store.Close()

	// Write at given offset.
	// offset must be a multiplier of 8, beacuse each uint64 will allocate 8 bytes.
	store.WriteAt(0, 10)  // mmap[0 : 8] = 10
	store.WriteAt(8, 20)  // mmap[8 : 16] = 20
	store.WriteAt(16, 34) // mmap[16 : 24] = 34

	// Read from a given offset
	num, _ := store.ReadAt(0)
	fmt.Println(num, "stored at offset 0")

	num, _ = store.ReadAt(16)
	fmt.Println(num, "stored at offset 16")

	// Remove file if you dont want to use later
	os.Remove(filename)
}
