# 🧠 Key-Value Store Roadmap (C/C++ + Go)

This roadmap guides you through building a **production-grade key-value store** using **C/C++ for the storage engine** and **Go for the API/monitoring layer**. It takes you from **scratch to intermediate** in system and backend programming.

---

## 📅 Phase 0: Environment Setup

### ✅ Goals
- Install C/C++ build tools: `gcc`, `g++`, `make`, `gdb`, `valgrind`
- Install Go toolchain: `go`, `protoc`, `cobra`, `prometheus`, `grpc`
- Editors: VSCode, CLion, or GoLand
- Learn how to compile and run basic C and Go programs

### 📘 Learn
- [Install Go](https://go.dev/doc/install)
- [GCC Toolchain on Linux](https://gcc.gnu.org/install/)
- [Using Make](https://www.gnu.org/software/make/manual/make.html)
- [Intro to gdb](https://sourceware.org/gdb/current/onlinedocs/gdb/)
- [Setting up Go project structure](https://golangprojectstructure.com/)

---

## 📅 Phase 1: Systems Fundamentals (Week 1–3)

### ✅ Learn
| Topic | C/C++ | Go |
|-------|-------|----|
| File I/O | `open`, `read`, `write`, `fopen`, `mmap` | `os.File`, `bufio`, `syscall.Mmap` |
| Memory Management | `malloc`, `free`, pointer arithmetic | slices, maps, garbage collection |
| Binary Encoding | `struct`, `fread/fwrite` | `encoding/binary`, `json`, `protobuf` |
| CLI Tooling | `argc`, `argv`, `getopt` | `cobra`, `flag` |
| Crash Consistency | `fsync`, `O_SYNC`, checksums | `defer`, `Flush()` on close |

### 🛠 Projects
- Build a CLI tool in C: `put <key> <value>` → append to a binary file
- Implement a struct-based binary serializer and deserializer
- Write a WAL (Write-Ahead Log) with CRC32 and replay logic
- Go project: map-backed in-memory store with JSON persistence

### 📘 Resources
- [Beej’s Guide to File I/O](https://beej.us/guide/bgc/)
- Linux System Programming – Robert Love (Ch. 2–4)
- [cstack's DB Tutorial](https://cstack.github.io/db_tutorial/)
- [Go by Example](https://gobyexample.com/)
- [Protobuf in Go](https://developers.google.com/protocol-buffers/docs/gotutorial)

---

## 📅 Phase 2: Core Engine - Index, WAL, SSTable (Week 4–6)

### ✅ Learn
| Component | What to Learn |
|----------|----------------|
| LSM Tree | MemTable (sorted map), WAL, SSTable format, compaction |
| B+ Tree (Alt.) | Internal/leaf nodes, disk-backed pages, split/merge |
| Paging | mmap + 4KB pages, fixed-size records |
| Serialization | Protobuf vs custom binary format |
| Concurrency | pthreads (mutexes), Go: goroutines, sync.Mutex |

### 🛠 Projects
- Implement WAL + MemTable (C)
- Build SSTable writer + memory-mapped reader (Go)
- Implement Mini LSM or B+ Tree
- Simulate a simple buffer pool with LRU page eviction

### 📘 Resources
- [CMU 15-445](https://15445.courses.cs.cmu.edu/fall2023/)
- [LSM Tree Paper (O'Neil)](https://www.cs.umb.edu/~poneil/lsmtree.pdf)
- [MIT B+Tree Lecture](https://ocw.mit.edu/courses/6-830/lectures/)
- [Bitcask Paper (Riak)](https://riak.com/assets/bitcask-intro.pdf)
- [Understanding mmap()](https://man7.org/linux/man-pages/man2/mmap.2.html)
- [Hash Tables in C (Tutorial)](https://github.com/jamesroutley/write-a-hash-table)

---

## 📅 Phase 3: gRPC API + Monitoring + Polish (Week 7–8)

### ✅ Learn
| Task | Tool |
|------|------|
| gRPC API | `grpc-go`, `.proto`, `context`, `net/http` |
| Metrics | Prometheus, `promhttp`, counters, histograms |
| Go <-> C | `cgo`, `C.CString`, `C.GoString`, memory management |
| CLI | `spf13/cobra`, `flag` |

### 🛠 Projects
- gRPC Server: `Put`, `Get`, `Delete` methods
- Expose Prometheus metrics: keys, mem usage, WAL size, ops/sec
- Implement graceful shutdown that flushes WAL and memtable
- Use `cgo` to wrap C engine functions and expose to Go
- Build CLI client to interact with gRPC service

### 📘 Resources
- [gRPC in Go](https://grpc.io/docs/languages/go/)
- [Prometheus Client for Go](https://prometheus.io/docs/guides/go-application/)
- [cgo Tutorial](https://www.ardanlabs.com/blog/2013/10/cgo-and-go-packages.html)
- [Golang Cobra CLI](https://github.com/spf13/cobra)

---

## 🧩 Final Projects (Pick One)

| Name | Description |
|------|-------------|
| RocksLite | LSM-based KV Store with WAL, SSTable, and compaction |
| MiniDocDB | Document DB with versioned blobs and query support |
| LiteDB | B+Tree-based embedded DB with binary storage |
| KVX | Go-only project: JSON-backed in-memory store |

---

## 🧠 What to Learn Per Language

### ✅ C / C++
- File I/O: `open`, `read`, `write`, `fopen`, `mmap`, `fsync`
- Struct packing, padding, and binary layout
- Manual memory management: `malloc`, `free`
- Tree structures: B+ Tree, Trie, custom hash table
- Threading with pthreads and mutex locks

### ✅ Go
- File operations with `os`, `bufio`, `encoding/binary`
- Slices, maps, goroutines, channels
- gRPC, Prometheus, cgo
- CLI tools with Cobra, graceful shutdowns
- Go module layout, build & testing

---

## 📚 Key Resources

### 📘 Books
- Linux System Programming – Robert Love
- The C Programming Language – Brian Kernighan & Dennis Ritchie
- Designing Data-Intensive Applications – Martin Kleppmann
- Advanced Programming in the UNIX Environment (APUE) – W. Richard Stevens
- Operating Systems: Three Easy Pieces (OSTEP) – Free [Online](http://pages.cs.wisc.edu/~remzi/OSTEP/)

### 🌐 Websites
- https://cstack.github.io/db_tutorial/
- https://build-your-own.org/database/
- https://beej.us/guide/bgnet/
- https://github.com/golang/protobuf
- https://github.com/prometheus/client_golang
- https://grpc.io/

- https://build-your-own.org/redis/
- https://beej.us/guide/bgnet/pdf/bgnet_usl_c_1.pdf

---

## 💡 Tips for C <-> Go Integration
- Keep performance-critical parts in C, interface in Go
- Use `C.CString()` and `C.GoString()` to convert strings
- Free memory manually allocated in C with `C.free()`
- NEVER pass Go pointers to C
- Use `defer` in Go to ensure cleanup (e.g. close WAL)

---

## ✅ Bonus Features
- AES-256 encryption for SSTables
- SSTable compaction and index merging
- WAL replication or Raft protocol integration
- Go-based admin dashboard (Grafana/Prometheus)
- On-disk persistent B+Tree with cache + LRU policy

---

## 📦 Tools You’ll Use
| Category | Tools |
|---------|-------|
| Build & Debug | `make`, `gdb`, `valgrind`, `lldb`, `cmake` |
| Go | `cobra`, `grpc`, `prometheus`, `cgo`, `pprof` |
| Serialization | `protobuf`, `encoding/binary`, JSON |
| Editors | VSCode, CLion, GoLand, Vim |

---

## ✅ Final Outcome
After completing this roadmap, you'll have:
- A fully working file-backed key-value store (WAL + SSTable)
- Clean C/C++ system code and Go API wrapper
- Real-world exposure to binary encoding, mmap, crash safety
- gRPC APIs, metrics, CLI tools, and documentation
- A portfolio-ready project you can deploy and extend


## GOLANG

| Topic           | Description                                     |
| --------------- | ----------------------------------------------- |
| File I/O        | `os.Open`, `os.Create`, `Read`, `Write`, `Seek` |
| Buffering       | `bufio.Reader`, `bufio.Writer`                  |
| Binary Encoding | `encoding/binary`, `struct layout`, endianness  |
| JSON            | `encoding/json` for dump/load memory stores     |
| Memory Mapping  | `syscall.Mmap` (for SSTables)                   |

Resources:
    [File I/O in Go](https://yourbasic.org/golang/reading-files/)
    [Binary encoding](https://golangdocs.com/binary-read-write)
    [mmap in Go](https://pkg.go.dev/golang.org/x/exp/mmap) (via syscall) 


| Topic      | Description                            |
| ---------- | -------------------------------------- |
| Goroutines | `go` keyword for concurrent execution  |
| Channels   | Message passing and synchronization    |
| Select     | Waiting on multiple channel operations |
| Mutexes    | `sync.Mutex`, shared memory            |
| WaitGroup  | Coordinating goroutines termination    |

Resources:
    [Go Concurrency Patterns](https://go.dev/blog/pipelines) (Go Blog)
    [Go’s sync package](https://pkg.go.dev/sync)

| Topic           | Description                                    |
| --------------- | ---------------------------------------------- |
| Protobuf        | `.proto` files, message types                  |
| gRPC Server     | `grpc.NewServer()`, services, request handling |
| gRPC Client     | Dialing, calling methods, context              |
| REST (Optional) | `net/http`, `mux`, `fiber`, JSON API           |
| Context         | `context.WithTimeout`, `context.Background`    |

Resources:
    [gRPC in Go Docs](https://grpc.io/docs/languages/go/)
    [Protocol Buffers Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
    [Go net/http](https://pkg.go.dev/net/http)
    [Go pprof Guide](https://go.dev/blog/pprof)

| Topic             | Description                                |
| ----------------- | ------------------------------------------ |
| Basic cgo         | `import "C"`, `C.CString`, `C.GoString`    |
| Memory management | `C.free()`, avoiding Go → C pointer leaks  |
| Linking           | Building `.so` shared libs, static linking |
| Errors            | Handle C error return values in Go idioms  |

Resources:
    [cgo docs](https://pkg.go.dev/cmd/cgo)
    [Ardan Labs cgo tutorial](https://www.ardanlabs.com/blog/2013/10/cgo-and-go-packages.html)
    [Golang + C Interop Tutorial](https://dev.to/salmanulfarzy/interfacing-go-and-c-using-cgo-19he)

## C/C++

| Topic             | Description                                          |
| ----------------- | ---------------------------------------------------- |
| Variables & Types | `int`, `char`, `float`, `struct`, `union`, `typedef` |
| Functions         | Parameters, return values, scope                     |
| Pointers          | Pointer arithmetic, double pointers, `NULL`          |
| Memory            | `malloc`, `calloc`, `free`, stack vs heap            |
| Arrays & Strings  | String manipulation, buffer overflow safety          |
| Header Files      | `#include`, declarations vs definitions              |

| Concept              | System APIs                                         |
| -------------------- | --------------------------------------------------- |
| File I/O             | `open()`, `read()`, `write()`, `lseek()`, `close()` |
| Standard I/O         | `fopen`, `fread`, `fwrite`, `fprintf`               |
| Struct Serialization | Padding, `#pragma pack`, binary layout              |
| mmap                 | `mmap()`, `munmap()`, `mprotect()`                  |
| File Permissions     | `mode_t`, `chmod`, `stat`                           |
| Endianness           | `htonl`, `ntohl`, portable storage format           |

Resources:
    📖 Linux System Programming – Robert Love (Ch. 2–4, 8)
    [Beej’s Guide to File I/O](https://beej.us/guide/bgc/)
    [Mortoray Binary Format Article](https://mortoray.com/2013/06/12/understanding-binary-files-in-c/)

Practice:
    Write binary key-value records to a file
    Use struct and fwrite() for record layout

| Structure       | Purpose                                |
| --------------- | -------------------------------------- |
| Linked List     | Simple page/record chains              |
| Hash Table      | Fast in-memory index                   |
| B+ Tree         | On-disk index for SSTable / table scan |
| Queue/Stack     | For LRU, WAL replay                    |
| Trie (optional) | Prefix search engine                   |

Resources:
    [GeeksforGeeks B+ Tree Tutorial](https://www.geeksforgeeks.org/b-tree-set-1-introduction-2/)
    [Write a Hash Table in C](https://github.com/jamesroutley/write-a-hash-table)

Practice:
    Build an in-memory B+ Tree with node split/merge
    Implement a memory-efficient hash table

| Topic             | What to Learn                                  |
| ----------------- | ---------------------------------------------- |
| WAL               | Append-only log for crash recovery             |
| Checksums         | Use CRC32 or SHA1 for data validation          |
| fsync             | Use to force disk flush                        |
| Journaling        | Write-before-commit consistency                |
| Atomic Operations | Use `rename`, `O_SYNC` to ensure atomic writes |

Resources:
    [Bitcask Paper](https://riak.com/assets/bitcask-intro.pdf) (used in Riak)
    [How WAL Works](https://en.wikipedia.org/wiki/Write-ahead_logging)
    [cstack WAL Code](https://cstack.github.io/db_tutorial/)

| Feature     | Why it Matters                            |
| ----------- | ----------------------------------------- |
| `mmap()`    | Memory-mapped files for high-speed I/O    |
| LRU Cache   | Reuse memory pages efficiently            |
| Paging      | Simulate 4KB page-based files like SQLite |
| Buffer Pool | Cache hot pages, evict cold pages         |

Resources:
    [Linux mmap Tutorial](https://man7.org/linux/man-pages/man2/mmap.2.html)
    [Memory Pager in C](https://cstack.github.io/db_tutorial/) (cstack)

| Tool        | Description                            |
| ----------- | -------------------------------------- |
| `make`      | Build system for modular C projects    |
| `gdb`       | Debugger for runtime tracing           |
| `valgrind`  | Detect memory leaks and pointer errors |
| Code Layout | Separate `.h` and `.c`, use modules    |
| Testing     | Write small testable units             |

Resources:
    [Makefile Tutorial](https://makefiletutorial.com/)
    [GDB Manual](https://sourceware.org/gdb/current/onlinedocs/gdb/)
    [Valgrind Quickstart](https://valgrind.org/docs/manual/quick-start.html)

| Concept               | How it works                              |
| --------------------- | ----------------------------------------- |
| Exporting C functions | Use `extern` and headers                  |
| cgo in Go             | `import "C"`, `//export MyFunc`           |
| Memory safety         | Use `C.CString`, `C.GoString`, `C.free()` |
| Compiling shared libs | `gcc -shared -fPIC` for `.so` files       |

Resources:
    [Go cgo Docs](https://pkg.go.dev/cmd/cgo)
    [cgo Shared Lib Guide](https://blog.kowalczyk.info/article/JyRZ/embedding-c-in-go.html)

Recommended Mini-Projects (C Only)
    Binary Record Writer – Write/read binary structs from file
    Simple WAL System – Append-only log with replay and checksum
    Buffer Pool Pager – Simulate 4KB page cache with LRU
    In-Memory B+ Tree – Implement insert, search, and split
    Shared Lib + Go – Export a C API and call it from Go using cgo