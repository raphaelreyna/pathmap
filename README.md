# pathmap

A Go package for creting nested map data structures from path strings.

## Example

```go
package main
import (
    "fmt"
    "github.com/raphaelreyna/pathmap"
)

func main() {
    m := pathmap.New()
    m.Set("a.b.c", 1)
    m.Set("a.e", 3)
    m.Delete("a.b")

    // We expect a.b.c to not be present
    // because a.b was deleted (and so were all its children)
    _, ok := m.Get("a.b.c")
    fmt.Println(ok)
    
    val, _ := m.Get("a.e")
    fmt.Println(val) // 3
    val := m["a"]["e"]
    fmt.Println(val) // 3
}
```