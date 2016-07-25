
## Componenta / Sqlx

Набросок мини пакета для конструктора запросов.

Sqlx используется совместно с пакетом database/sql

Пакет поддерживает СУБД PostgreSql и MySql

```go
package main

import (
    "database/sql"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

func main() {
    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

	sqlx.Driver("mysql")

	// SELECT `id`, `name` FROM `users` WHERE `age` > ? ORDER BY `created_at` DESC LIMIT ? OFFSET ?;
	query := sqlx.Table("users").
		Where("age", ">", 21).
		OrderBy("created_at", "DESC").
		Limit(15).
		Offset(30).
		Select("id", "name")

    rows, err := db.Query(query.Sql(), query.Data()...)

	// ...
	// ...
}
```

Больше примеров в Doc файле.
