
## Componenta / Hash

Обертка над SHA256... остальные позже

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/hash"
)

func main() {
    h := hash.New("String")
    sum := h.Sum([]byte("String"))

    fmt.Println(sum)
}
```