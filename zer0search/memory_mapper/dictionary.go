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

type Dictionary struct {
	file   *os.File
	mmap   gommap.MMap // mmap
	len    uint64      // current size
	closed bool        // flag to check if the directory.index is closed
}

func NewDictionary() (*Dictionary, error) {
	os.Remove(filepath.Join(utils.Path, dictIndexFile))
	dict := &Dictionary{}
	file, err := os.OpenFile(filepath.Join(utils.Path, dictIndexFile), os.O_CREATE|os.O_RDWR, 0644)
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

// search in dictionary.index
// [wordHash][postingOffset][postingLength]
// [uint64][uint64][uint64]
// [8][8][8] -> 24bytes
// returns offsetOfWord, postingOffset, postingLen
func (d *Dictionary) Search(wordHash uint64) (bool, uint64, uint64, uint64, error) {
	if d.closed {
		return false, uint64(0), uint64(0), uint64(0), errors.New("dictionary.index file is closed")
	}
	for j := uint64(0); j+dictEntrySize <= d.len; j += dictEntrySize {
		offset := j
		storedHash := encoder.Uint64(d.mmap[offset : offset+byteSize])
		offset += byteSize
		postingOffset := encoder.Uint64(d.mmap[offset : offset+byteSize])
		offset += byteSize
		positionLength := encoder.Uint64(d.mmap[offset : offset+byteSize])
		if storedHash == wordHash {
			return true, uint64(j), postingOffset, positionLength, nil
		}
	}
	return false, uint64(0), uint64(0), uint64(0), nil
}

// append in dictionary.index
func (d *Dictionary) Append(hash, postingOffset, postingLen uint64) error {
	if d.closed {
		return errors.New("file is closed")
	}
	if d.IsFilled(dictEntrySize) {
		return errors.New("max filesize reached")
	}
	offset := d.len
	encoder.PutUint64(d.mmap[offset:offset+byteSize], hash)
	offset += byteSize
	encoder.PutUint64(d.mmap[offset:offset+byteSize], postingOffset)
	offset += byteSize
	encoder.PutUint64(d.mmap[offset:offset+byteSize], postingLen)
	offset += byteSize
	d.len += dictEntrySize
	return nil
}

// update postingOffset and postingLen stored at offset
func (d *Dictionary) Update(offset, postingOffset, postingLen uint64) error {
	if d.closed {
		return errors.New("file is closed")
	}
	if offset+dictEntrySize > d.len {
		return errors.New("[error] : SGMNT_FLT")
	}
	initialOffset := offset + byteSize
	encoder.PutUint64(d.mmap[initialOffset:initialOffset+byteSize], postingOffset)
	initialOffset += byteSize
	encoder.PutUint64(d.mmap[initialOffset:initialOffset+byteSize], postingLen)
	return nil
}

// check if there is enough space with size
func (d *Dictionary) IsFilled(size uint64) bool {
	return d.len+size > MaxFileSize
}

// close the dictionary.index
func (d *Dictionary) Close() error {
	slog.Info("[INSIDE] dictionary.go -> Close()")
	if d.closed {
		return errors.New("file is closed")
	}
	if err := d.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := d.file.Sync(); err != nil {
		return err
	}
	if err := d.file.Truncate(int64(d.len)); err != nil {
		return err
	}
	if err := d.file.Sync(); err != nil {
		return err
	}
	if err := d.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	if err := d.file.Close(); err != nil {
		return err
	}
	d.closed = true
	d.mmap = nil
	d.file = nil
	return nil
}

// debug information
func (d *Dictionary) Debug() {
	for j := uint64(0); j+dictEntrySize <= d.len; j += dictEntrySize {
		offset := j
		storedHash := encoder.Uint64(d.mmap[offset : offset+byteSize])
		offset += byteSize
		postingOffset := encoder.Uint64(d.mmap[offset : offset+byteSize])
		offset += byteSize
		positionLength := encoder.Uint64(d.mmap[offset : offset+byteSize])
		fmt.Println("offset: ", j, "shoredHash: ", storedHash, "postingOffset: ", postingOffset, "postingLen: ", positionLength)
	}
}
