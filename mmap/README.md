# MMAP

## What is mmap ?
Instead of reading a file into RAM and keeping a copy, mmap lets you treat the file like part of your program’s memory, leveraging the OS’s paging mechanism: only the accessed pages are loaded into physical memory. Only usable in UNIX based OS. mmap is just a slice of bytes.

## Advantages:
- https://www.tencentcloud.com/techpedia/106444
- Saves ram
- Loads a specific parts
- No read(), write() syscall, directly accessing memory allowing fast read and writing. check  https://learningdaily.dev/reading-and-writing-files-using-memory-mapped-i-o-220fa802aa1c
- Update data in memory and it's saved to disk automatically by using flushing
- Multiple process can map the same file and share data

## Application:
    - huge files (logs, indexes, databases)
    - shared memory communication

## Implementation Idea:
    - autocomplete from millions of keywords
    - persistent key-value store (modern databases uses this mechanism)
    - process a 10gb log file
    - memory allocator that uses disk-based memory

## Resources:
    - https://en.wikipedia.org/wiki/Memory-mapped_file
    - https://www.tutorialspoint.com/unix_system_calls/mmap.htm
    - https://kuafu1994.github.io/MoreOnMemory/shared-memory.html
    - https://blog.minhazav.dev/memory-sharing-in-linux/
    - https://medium.com/%40jyjimmylee/how-does-memory-mapping-mmap-work-c8a6a550ba0d
    - https://stackoverflow.com/questions/258091/when-should-i-use-mmap-for-file-access
    - https://medium.com/cosmos-code/mmap-vs-read-a-performance-comparison-for-efficient-file-access-3e5337bd1e25

## Install package
    go get github.com/tysonmote/gommap

## Memory usages of os.File
For a 10GB file, os.File will store some information about that file (some bytes), when we try to read from this file then it will use memory. 