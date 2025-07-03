package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

// This program implements a persistent, file-backed Trie using memory-mapped files.
// Each node is stored as a fixed-size binary structure within the file and is accessed using file offsets.
// This avoids storing the Trie in RAM, enabling scalable, disk-based data structures.

var encoder = binary.BigEndian // consistent byte order(Big Endian) for binary encoding

// Node represents a single trie node.
// - Exits: 1 if this node marks the end of a word, 0 otherwise
// - Child: array of offsets to child nodes for each lowercase letter ('a' to 'z')
type Node struct {
	Exits uint8      // 1byte
	Child [26]uint64 // 26 * 8 = 208 bytes
}

// memory-mapped trie structure
type Trie struct {
	file       *os.File
	mmap       gommap.MMap
	nextOffset uint64 // offset to store the next node
	entrySize  uint64 // size of each node 209 byte
	maxSize    uint64
}

// NewTrie creates a new memory-mapped trie with the given filename and maximum size.
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

	// Expand file to max size
	if err := file.Truncate(int64(maxSize)); err != nil {
		file.Close()
		return nil
	}

	// mmap the file with read/write access, shared among processes
	mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil
	}

	return &Trie{
		file:       file,
		mmap:       mmap,
		nextOffset: uint64(info.Size()),
		entrySize:  uint64(binary.Size(Node{})),
		maxSize:    uint64(maxSize),
	}
}

// ReadAt reads a Node struct from the mmap'd file at the given offset.
func (t *Trie) ReadAt(offset uint64) Node {
	node := Node{}
	node.Exits = t.mmap[offset]
	offset++
	for i := 0; i < 26; i++ {
		node.Child[i] = encoder.Uint64(t.mmap[offset : offset+8])
		offset += 8
	}
	return node
}

// WriteAt writes a Node struct to the mmap'd file at the given offset.
func (t *Trie) WriteAt(offset uint64, node Node) {
	t.mmap[offset] = node.Exits
	offset++
	for i := 0; i < 26; i++ {
		encoder.PutUint64(t.mmap[offset:offset+8], node.Child[i])
		offset += 8
	}
}

// NextOffset returns the next available offset for writing a new node.
// It ensures that we do not exceed the maximum file size.
func (t *Trie) NextOffset() (uint64, error) {
	curOffset := t.nextOffset
	t.nextOffset += t.entrySize
	if t.nextOffset > t.maxSize {
		return 0, fmt.Errorf("overflow")
	}
	return curOffset, nil
}

// Insert inserts a lowercase word into the trie starting from the given root offset.
// If the trie does not have a node for the next char, it allocates new nodes using NextOffset.
// Here is a catch, if the file reaches it's max size while storing ?
func (t *Trie) Insert(word string, rootOffset uint64) error {
	curOffset := rootOffset
	curNode := t.ReadAt(curOffset)
	for _, c := range word {
		index := c - 'a'
		if curNode.Child[index] == 0 {
			newOffset, err := t.NextOffset() // create a new node
			if err != nil {
				return err
			}
			curNode.Child[index] = newOffset
			t.WriteAt(curOffset, curNode)
		}
		curOffset = curNode.Child[index]
		curNode = t.ReadAt(curOffset)
	}
	curNode.Exits = 1 // mark as ending node
	t.WriteAt(curOffset, curNode)
	return nil
}

// Search checks whether a given lowercase word exists in the Trie.
// It starts traversal from the given rootOffset and follows each character.
// Returns true if the word is found and its terminal node has Exits = 1.
func (t *Trie) Search(word string, rootOffset uint64) bool {
	curOffset := rootOffset
	curNode := t.ReadAt(curOffset)
	for _, c := range word {
		index := c - 'a'
		if curNode.Child[index] == 0 { // node doesn't exists
			return false
		}
		curOffset = curNode.Child[index]
		curNode = t.ReadAt(curOffset)
	}
	return curNode.Exits == uint8(1)
}

// Close ensures all changes are flushed to disk, unmaps memory, and closes the file.
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

	// Create a new trie with 10MB trie.index file
	trie := NewTrie("trie.index", 10*1024*1024)
	defer trie.Close()

	// Root node of the trie
	rootOfset, err := trie.NextOffset()
	if err != nil {
		panic(err)
	}

	fmt.Println(rootOfset)

	// Insert word
	trie.Insert("labib", rootOfset)
	trie.Insert("labix", rootOfset)

	// Search word
	fmt.Println(trie.Search("labib", rootOfset))
	fmt.Println(trie.Search("lab", rootOfset))
	fmt.Println(trie.Search("labiba", rootOfset))
	fmt.Println(trie.Search("a", rootOfset))

	// Remove if you don't want to use the trie later
	os.Remove("trie.index")
}
