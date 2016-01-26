
## Componenta / Pass

Хеширование паролей по алгоритму Blowfish.

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/pass"
    "log"
)

func main() {
    p := pass.New("String")
    h, err := p.Hash("Password")
    
    if err != nil {
        log.Fatalln("Pass:", err)
    }
    
    if p.Compare(h, "Password") {
        fmt.Println("Compare!")
    }
}
```