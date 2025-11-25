# fio package - clone of fmt package
---
Supports only 

| Concrete Type  |
|----------------|
|       Int      |
|     String     |
|     Struct     |

---
Usage :
```go
import "fio"

type user struct {
    Name string
    Age  int
}

var s string = "Hello"
var n int = 12
var bob user = user {
    Name : "Bob",
    Age : 12,
}

fio.Write(s, n, bob)
fio.Fwrite(fio.Out, s, n, bob)

fio.Read(&s, &n)
```

```go
var b fio.ByteBuilder

b.WriteString("hello")
b.WriteBytes([]byte(" world"))
b.WriteByte('\n')

fio.Write(b.String())
fio.Write(b.Byte())
```