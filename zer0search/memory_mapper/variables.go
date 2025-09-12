package memorymapper

import "encoding/binary"

var (
	dictIndexFile           = "/memory_mapper/dictionary.index"
	postingIndexFile        = "/memory_mapper/posting.index"
	byteSize         uint64 = 8
	dictEntrySize    uint64 = 24       // [hash][offset][postingLen]
	MaxFileSize      uint64 = 10485760 // 10Mb

	encoder = binary.BigEndian
)
