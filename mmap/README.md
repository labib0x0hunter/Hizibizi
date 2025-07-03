# MMAP

## What is mmap ?
Instead of reading a file into RAM and keeping a copy, mmap lets you treat the file like part of your program’s memory (virtual memory), leveraging the OS’s paging mechanism: only the accessed pages are loaded into physical memory. Only usable in UNIX based OS. mmap is just a slice of bytes.
[Learn more](https://en.wikipedia.org/wiki/Memory-mapped_file)

## Advantages:
- [Advantages](https://www.tencentcloud.com/techpedia/106444)
- Saves ram
- Loads a specific parts
- No read(), write() syscall, directly accessing memory allowing fast read and writing. [Check](https://learningdaily.dev/reading-and-writing-files-using-memory-mapped-i-o-220fa802aa1c)
- Update data in memory and it's saved to disk automatically by using flushing
- Multiple process can map the same file and share data

## Application:
- Processing huge files (logs, indexes, databases)
- Shared memory communication
- Persistent key-value store (modern databases uses this mechanism) [Check](https://brunocalza.me/2021/01/18/but-how-exactly-databases-use-mmap)

## Resources:
- https://www.tutorialspoint.com/unix_system_calls/mmap.htm
- https://kuafu1994.github.io/MoreOnMemory/shared-memory.html
- https://blog.minhazav.dev/memory-sharing-in-linux/
- https://medium.com/%40jyjimmylee/how-does-memory-mapping-mmap-work-c8a6a550ba0d
- https://stackoverflow.com/questions/258091/when-should-i-use-mmap-for-file-access
- https://medium.com/cosmos-code/mmap-vs-read-a-performance-comparison-for-efficient-file-access-3e5337bd1e25
- [Byte Order](https://betterexplained.com/articles/understanding-big-and-little-endian-byte-order/)

## Install package
    go get github.com/tysonmote/gommap

## Memory usages of os.File
For a 10GB file, os.File will store some information about that file (some bytes), when we try to read from this file then it will use memory.

## Use cases
```go
mmap, err := gommap.Map(file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)

// mmap is a []byte
// So, mapping a 15-byte sized file, mmap[offset] is one byte.
// Visualize mmap as, a slice of bytes with len 15 and offset = index


// store a string
offset := 0
s := []byte("hello")
copy(mmap[offset: offset + len(s)], s)
offset += len(s)    // advance offset by 5

// store a int
num := uint32(34)   // 4bytes
binary.BigEndian.PutUint32(mmap[offset: offset + 4], num)   // encode num into []byte

// retrive string
data := mmap[0: len(s)] // we stored string at offset 0
fmt.Println(string(data))

// retrive int
storedNum := binary.BigEndian.Uint32(mmap[offset: offset + 4])
fmt.Println(storedNum)
```

