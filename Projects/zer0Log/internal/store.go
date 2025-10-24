package internal

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	File *os.File      // file to store log
	mu   sync.Mutex    // threat-safely writes
	buf  *bufio.Writer // buffer writer
	size uint64        // current size of file
}

func NewStore(file *os.File) (*store, error) {
	info, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(info.Size())
	return &store{
		File: file,
		buf:  bufio.NewWriter(file),
		size: size,
	}, nil
}

func (s *store) Append(data []byte) (uint64, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos := s.size
	if err := binary.Write(s.buf, enc, uint64(len(data))); err != nil {
		return 0, 0, err
	}

	w, err := s.buf.Write(data)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	data := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(data, int64(pos + lenWidth)); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}