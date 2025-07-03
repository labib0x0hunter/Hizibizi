package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

var encoder = binary.BigEndian

type Node struct {
	Exits uint8
	Child [26]uint64
}

type Trie struct {
	file       *os.File
	mmap       gommap.MMap
	nextOffset uint64
	entrySize  uint64 // 209 + 1 byte
	maxSize    uint64
}

func NewTrie(filename string, maxSize uint64) *Trie {
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
	if err := file.Truncate(int64(maxSize)); err != nil {
		file.Close()
		return nil
	}

	// mmap the file
	mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil
	}

	return &Trie{
		file:       file,
		mmap:       mmap,
		nextOffset: uint64(info.Size()),
		entrySize:  uint64(binary.Size(Node{})) + uint64(1),
		maxSize:    uint64(maxSize),
	}
}

func (t *Trie) ReadAt(offset uint64) Node {
	node := Node{}
	node.Exits = uint8(encoder.Uint16(t.mmap[offset : offset+2]))
	offset += 2
	for i := 0; i < 26; i++ {
		node.Child[i] = encoder.Uint64(t.mmap[offset : offset+8])
		offset += 8
	}
	return node
}

func (t *Trie) WriteAt(offset uint64, node Node) {
	encoder.PutUint16(t.mmap[offset:offset+2], uint16(node.Exits))
	offset += 2
	for i := 0; i < 26; i++ {
		encoder.PutUint64(t.mmap[offset:offset+8], node.Child[i])
		offset += 8
	}
}

func (t *Trie) NextOffset() (uint64, error) {
	curOffset := t.nextOffset
	t.nextOffset += t.entrySize
	if t.nextOffset > t.maxSize {
		return 0, fmt.Errorf("overflow")
	}
	return curOffset, nil
}

// Always insert lowercase word
func (t *Trie) Insert(word string, rootOffset uint64) error {
	curOffset := rootOffset
	curNode := t.ReadAt(curOffset)
	for _, c := range word {
		index := c - 'a'
		if curNode.Child[index] == 0 {
			newOffset, err := t.NextOffset()
			if err != nil {
				return err
			}
			curNode.Child[index] = newOffset
			t.WriteAt(curOffset, curNode)
		}
		curOffset = curNode.Child[index]
		curNode = t.ReadAt(curOffset)
	}
	curNode.Exits = 1
	t.WriteAt(curOffset, curNode)
	fmt.Println(curOffset)
	return nil
}

func (t *Trie) Close() error {
	if err := t.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := t.file.Sync(); err != nil {
		return err
	}
	if err := t.file.Truncate(int64(t.nextOffset)); err != nil {
		return err
	}
	if err := t.file.Sync(); err != nil {
		return err
	}
	if err := t.mmap.UnsafeUnmap(); err != nil {
		return err
	}
	if err := t.file.Close(); err != nil {
		return err
	}
	return nil
}

func main() {

	trie := NewTrie("trie.index", 10*1024*1024)
	defer trie.Close()

	rootOfset, err := trie.NextOffset()
	if err != nil {
		panic(err)
	}

	fmt.Println(rootOfset)

	trie.Insert("labib", rootOfset)
	trie.Insert("labix", rootOfset)

	os.Remove("trie.index")
}
