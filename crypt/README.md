
## Componenta / Crypt

Пакет для шифрование с использованием секретного ключа. Для шифрования используется алгоритм Advanced Encryption Standard (AES), также известный как Rijndae, принятый в качестве стандарта шифрования в США. 

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/crypt"
    "log"
)

func main() {
    c, err := crypt.Encrypt("String", "Secret_Key")

    if err != nil {
        log.Fatalln("Encrypt:", err)
    }
    
    fmt.Println(c)
    
    s, err := crypt.Decrypt(c, "Secret_Key")

    if err != nil {
        log.Fatalln("Decrypt:", err)
    }
    
    fmt.Println(s)
}
```