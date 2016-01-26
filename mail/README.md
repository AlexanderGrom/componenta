
## Componenta / Mail

Набросок мини пакета для отправки почты.

```go
package main

import (
    "github.com/AlexanderGrom/componenta/mail"
)

func main() {
    m := mail.New()
    m.To("jack@example.com", "Jack")
    m.From("bob@example.com", "Bob")
    m.Subject("Hello")
    m.Text("Hello World")
    m.Send()
}
```