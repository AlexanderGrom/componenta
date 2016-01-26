
## Componenta / DBx

Улучшаем работу со стандартным пакетом `database/sql`.
По сути нас не устраивает только способ выборки, а именно скан в переменные. Поэтому мы расширим тип `sql.Rows` и добавим более высокоуровневый тип `Fetch`.

#### Что имеем сейчас:

**Скан в структуру**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/dbx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type People struct {
    Id    int
    Name  string
    Age   int
}

var (
    people People
)

func main() {
    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    rows, err := db.Query("SELECT people_id, people_name, people_age FROM peoples WHERE people_id = 123 LIMIT 1")

    if err != nil {
        log.Fatalln("DB Query:", err)
    }
    
    // Оборачиваем sql.Rows
    fetch := dbx.NewFetch(rows)
    
    // Сканируем результат в структуру
    err = fetch.ScanStruct(&people)
    
    if err != nil {
        log.Fatalln("DB Fetch:", err)
    }

    fmt.Printf("People name: %s\n", people.Name)

    db.Close()
}
```

**Скан в срез структур**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/dbx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type People struct {
    Id    int
    Name  string
    Age   int
}

var (
    peoples []People
)

func main() {
    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    rows, err := db.Query("SELECT people_id, people_name, people_age FROM peoples ORDER BY people_add_date DESC LIMIT 10")

    if err != nil {
        log.Fatalln("DB Query:", err)
    }
    
    // Оборачиваем sql.Rows
    fetch := dbx.NewFetch(rows)
    
    // Сканируем результат в срез структур
    err = fetch.ScanSlice(&peoples)

    if err != nil {
        log.Fatalln("DB Fetch:", err)
    }
    
    for _, pls := range peoples {
        fmt.Printf("People #%d, %s, %d\n", pls.Id, pls.Name, pls.Age)
    }

    db.Close()
}
```

**Скан в переменные**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/dbx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

var (
    name string
    age  int
)

func main() {
    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    rows, err := db.Query("SELECT people_id, people_name, people_age FROM peoples WHERE people_id = 123 LIMIT 1")

    if err != nil {
        log.Fatalln("DB Query:", err)
    }
    
    // Оборачиваем sql.Rows
    fetch := dbx.NewFetch(rows)
    
    // Сканируем результат в переменные
    err = fetch.ScanVars(&name, &age)
    
    if err != nil {
        log.Fatalln("DB Fetch:", err)
    }

    fmt.Printf("People: %s, %d\n", name, age)

    db.Close()
}
```

**Больший контроль над выборкой**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/dbx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type People struct {
    Id    int
    Name  string
    Age   int
}

var (
    people People
)

func main() {
    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    rows, err := db.Query("SELECT people_id, people_name, people_age FROM peoples WHERE people_id = 123 LIMIT 1")

    if err != nil {
        log.Fatalln("DB Query:", err)
    }
    
    // Оборачиваем sql.Rows
    rs := dbx.NewRows(rows)
    
    // Сканируем
    for rs.Next() {
        rs.ScanStruct(&people)
        
        if err != nil {
            log.Fatalln("DB Scan:", err)
        }
        
        fmt.Printf("People: %s, %d\n", people.Name, people.Age)
    }

    db.Close()
}
```