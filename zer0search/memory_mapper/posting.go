package memorymapper

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"searchengine/utils"

	"github.com/tysonmote/gommap"
)

type Posting struct {
	file   *os.File
	mmap   gommap.MMap // mmap
	len    uint64      // current size
	closed bool        // flag to check if the posting.index is closed
}

func NewPosting() (*Posting, error) {
	os.Remove(filepath.Join(utils.Path, postingIndexFile))
	slog.Info(" [NewPosting] Path -> " + filepath.Join(utils.Path, postingIndexFile))
	dict := &Posting{}
	file, err := os.OpenFile(filepath.Join(utils.Path, postingIndexFile), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	dict.closed = false
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	dict.len = uint64(size)
	if size > int64(MaxFileSize) {
		return nil, errors.New("max file size reached")
	}
	// allocate capacity
	if err := file.Truncate(int64(MaxFileSize)); err != nil {
		return nil, err
	}
	dict.file = file
	mmap, err := gommap.Map(
		dict.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	)
	if err != nil {
		file.Close()
		return nil, err
	}
	dict.mmap = mmap
	return dict, nil
}

// search in posting.index
// [len][slice]
// [uint64][uint64]
// [8][8 * len]
// returns offset, docIds
func (p *Posting) Search(offset uint64, len uint64) ([]uint64, error) {
	if p.closed {
		return []uint64{}, errors.New("posting.index file is closed")
	}
	storedLen := encoder.Uint64(p.mmap[offset : offset+byteSize])
	if storedLen != len {
		return []uint64{}, errors.New("length size didn't match, maybe stored a wrong offset")
	}
	totalByte := (len * byteSize) + byteSize // lenSize for len, (len * lenSize) for slice
	if offset+totalByte > MaxFileSize {
		return []uint64{}, errors.New(" [INFO]posting.index does not have enough space")
	}
	docIds := make([]uint64, 0, len)
	for i := 0; i < int(len); i++ {
		offset += byteSize
		x := encoder.Uint64(p.mmap[offset : offset+byteSize])
		docIds = append(docIds, x)
	}
	return docIds, nil
}

// append in posting.index
func (p *Posting) Append(docId uint64, sizeInclude bool) (uint64, error) {
	if p.closed {
		return 0, errors.New("posting.index file is closed")
	}
	initialOffset := p.len
	totalByte := byteSize
	if sizeInclude {
		totalByte = byteSize + byteSize
	}
	if p.IsFilled(totalByte) {
		return 0, errors.New("posting.index file is reached to maximum size")
	}
	offset := initialOffset
	if sizeInclude {
		encoder.PutUint64(p.mmap[offset:offset+byteSize], uint64(1))
		offset += byteSize
	}
	encoder.PutUint64(p.mmap[offset:offset+byteSize], docId)
	p.len += totalByte
	return initialOffset, nil
}

// Meaning we have to append some docId in the slice located at offset
// But thats a probem ...
// So we will copy the slice to the end and append the docId
func (p *Posting) Update(offset, size, docId uint64) (uint64, error) {
	if p.closed {
		return 0, errors.New("file is closed")
	}
	initialOffset := p.len
	totalByte := byteSize + size*byteSize
	if p.IsFilled(totalByte + byteSize) { // mmap[offset: offset + totalbyte] -> for old slice, lenSize for new docId
		return 0, errors.New("file is reached to maximum size")
	}
	copy(p.mmap[initialOffset:initialOffset+totalByte], p.mmap[offset:offset+totalByte])
	encoder.PutUint64(p.mmap[initialOffset:initialOffset+byteSize], size+uint64(1))
	p.len += totalByte
	if _, err := p.Append(docId, false); err != nil {
		return 0, err
	}
	return initialOffset, nil
}

// check if there is enough space with size
func (p *Posting) IsFilled(size uint64) bool {
	return p.len+size > MaxFileSize
}

// close the posting.index
func (p *Posting) Close() error {
	slog.Info("[INSIDE] posting.go -> Close()")
	if p.closed {
		return errors.New("file is closed")
	}
	if err := p.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := p.file.Sync(); err != nil {
		return err
	}
	if err := p.file.Truncate(int64(p.len)); err != nil {
		return err
	}
	if err := p.file.Sync(); err != nil {
		return err
	}
	if err := p.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	if err := p.file.Close(); err != nil {
		return err
	}
	p.closed = true
	p.mmap = nil
	p.file = nil
	return nil
}

// Print slice for debug
func (p *Posting) Print(offset uint64) uint64 {
	sz := encoder.Uint64(p.mmap[offset : offset+byteSize])
	fmt.Println(" [debug] first slice size: ", sz)
	fmt.Print(" [debug] first slice size: ")
	for i := 0; i < int(sz); i++ {
		offset += byteSize
		x := encoder.Uint64(p.mmap[offset : offset+byteSize])
		fmt.Print(x, " ")
	}
	fmt.Println()
	return sz
}

// debug information
// customize to your need
func (p *Posting) Debug(offset uint64) {
	fmt.Println(" [debug] bytes written: ", p.len)
	fmt.Println(" [debug] bytes : ", p.len/8)

	sz := p.Print(uint64(offset))

	offset += byteSize
	offset += (byteSize * sz)

	fmt.Println(offset)
}

func (p *Posting) Len() uint64 {
	return p.len
}
