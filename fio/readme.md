# fio package - clone of fmt package
---
Supports only int, string, struct for minimalist

| Verb | Type   |
|------|--------|
| %s   | String |
| %d   | Int    |
| %S   | Struct |

---
Usage :
```go
import fio

type user struct {
    Name string
    Age  int
}

var s string = "Hello"
var n int = 12
var bob user = user{
    Name : "Bob",
    Age : 12,
}

fio.Write("%s %d %S\n", s, n, bob)
fio.Fwrite(fio.Out, "%s %d %S\n", s, n, bob)
```