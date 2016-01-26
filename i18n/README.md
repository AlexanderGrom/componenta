
## Componenta / i18n

Парсинг простых i18n файлов.

```ini
data.timezone = "Москва"
data.month = "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"
```

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/i18n"
)

func main() {
    lang, _ := i18n.Use("./i18n/ru/system.lang")
    fmt.Println(lang.GetString("data.timezone"))
    fmt.Println(lang.GetStrings("data.month"))
}
```