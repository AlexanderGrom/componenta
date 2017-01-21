
## Componenta / env

Установка переменных окружения из файла.

```ini
# Настройки приложения
APP_HOME = "/var/app"
APP_HOST = "localhost"
APP_PORT = "1448"
APP_URL  = "http://${APP_HOST}:${APP_PORT}" # http://localhost:1448
```

```go
package main

import (
    "fmt"
    "os"
    "github.com/AlexanderGrom/componenta/env"
)

func main() {
    env.Load("${HOME}/.env")

    fmt.Println(os.Getenv("APP_HOME"))
    fmt.Println(os.Getenv("APP_HOST"))
    fmt.Println(os.Getenv("APP_PORT"))
    fmt.Println(os.Getenv("APP_URL"))
}
```
