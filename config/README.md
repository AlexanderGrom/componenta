
## Componenta / Config

Парсинг простых конфигуровочных файлов.

```ini
key1 = value     # Ключ значение (default)
key2 = "value2"  # Строковые значения

key3 = 11   # Целочисленные значения
key4 = 3.14 # Вещественные значения
key5 = true # Булевые значения ("true", "1", "on", "yes")

key5 = 1,2,3,4,5,6,7,8,9 # Список целочисленных значений

# Обычные
# Комментарии...

key6 = true, false # Список булевых значений. 
                   # Могут быть:
                   # "true", "1", "on", "yes",
                   # "false", "0", "off", "no"
                   
key7 = "foo", "bar"  # Список строковых зачений

key8    =       # Пустое значение
key9            # Пустое значение
key10   =   ""  # Пустое значение
```

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/config"
)

func main() {
    cfg, _ := config.Use("./config/database.conf")

    fmt.Println(cfg.GetString("key1"))
    fmt.Println(cfg.GetBool("key5"))
    fmt.Println(cfg.GetInts("key5"))
    fmt.Println(cfg.GetStrings("key7"))
}
```