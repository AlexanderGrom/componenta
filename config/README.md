
## Componenta / Config

Парсинг простых конфигуровочных файлов.

```ini
key1 = value     # Ключ значение (default)
key2 = "value2"  # Строковые значения

key3 = 11   # Целочисленные значения
key4 = 3.14 # Вещественные значения
key5 = true # Булевые значения ("true", "1", "on", "yes")

key6 = 1,2,3,4,"5",6,7,8,9 # Список целочисленных значений

# Обычные
# Комментарии...

key7 = true, false # Список булевых значений.
                   # Могут быть:
                   # "true", "1", "on", "yes",
                   # "false", "0", "off", "no"

key8 = "foo" , "bar"  # Список строковых зачений

key9    =       # Пустое значение
key10           # Пустое значение
key11   =   ""  # Пустое значение

key12 = 1, 2, 3, 4, # Списки можно переносить на новую строку
        5, 6, 7, 8  # И оставлять комментарии...

key13 = "value
         value
         value" # Значения в кавычках могут состоять из нескольких строк

key14 = "value[\"key\"]" # Кавычки в строках можно экранировать

key15 = ${TESTVARNAME} | "value" # Если значение переменной окружения не пустое место, иначе дефолтное значение

key16 = ${EMPTYVARNAME} | value

key17 = 1,2,3, ${EMPTYVARNAME} | 4, 5, 6

key18 = ${EMPTYVARNAMEONE} | ${EMPTYVARNAMETWO} | value

key19 = | value

key20 = "${TESTVARNAME}/src" # Переменные могут быть в строках

key21 = "${TESTVARNAME/src"

key22 = /var/${TESTVARNAME}/src # Или в строках без кавычек

key23 = \${TESTVARNAME} # Экранирование
```

```go
package main

import (
    "fmt"
    "github.com/AlexanderGrom/componenta/config"
)

func main() {
    cfg, _ := config.Use("${HOME}/config/database.conf")

    fmt.Println(cfg.GetString("key1"))
    fmt.Println(cfg.GetBool("key5"))
    fmt.Println(cfg.GetInts("key5"))
    fmt.Println(cfg.GetStrings("key7"))
}
```
