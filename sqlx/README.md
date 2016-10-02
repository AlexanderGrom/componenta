
## Componenta / Sqlx

Мини пакет конструктора запросов.  
Пакет поддерживает СУБД PostgreSql, MySql и Sqlite

```go
package main

import (
    "database/sql"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    db, err := sql.Open("postgres", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

	sqlx.Driver("pgsql")

	// SELECT "id", "name" FROM "users" WHERE "age" > $1 ORDER BY "created_at" DESC LIMIT $2;
	query := sqlx.Table("users").
        Select("id", "name").
		Where("age", ">", 21).
		OrderBy("created_at", "DESC").
		Limit(10)

    rows, err := db.Query(query.Sql(), query.Data()...)

	// ...
	// ...

    db.Close()
}
```

**Ображение к талблице по её имени**
```go
// SELECT * FROM "users"
sql := sqlx.Table("users").Sql()
```

**Вложенный запрос в From**
```go
// SELECT * FROM (SELECT "group_id", MAX("created_at") as "lastdate" FROM "users" GROUP BY "group_id") as "users" ORDER BY "lastdate" DESC
sql := sqlx.Table(func(builder *sqlx.Builder) {
    builder.Select("group_id").From("users").GroupBy("group_id").Max("created_at", "lastdate")
}).OrderBy("lastdate", "DESC").Sql()
```

**Вложенный запрос в From используя разные построители**
```go
// SELECT * FROM (SELECT "group_id", MAX("created_at") as "lastdate" FROM "users" GROUP BY "group_id") as "users" ORDER BY "lastdate" DESC
sub := sqlx.Table("users").
    Select("group_id").
    GroupBy("group_id").
    Max("created_at", "lastdate")
sql := sqlx.Table(sqlx.Raw("( "+subque.Sql()+" ) as users", subque.Data()...)).
    OrderBy("lastdate", "DESC").
    Sql()
```

**Выборка Select**
```go
// SELECT "id", "name", "age" FROM "users"
sql := sqlx.Table("users").Select("id", "name", "age").Sql()
```

**Комбинированный Select**
```go
// SELECT "id", "name", "age", COUNT(*) as count FROM "users"
sql := sqlx.Table("users").
    Select("id", "name").Select("age").
    SelectRaw("COUNT(*) as count").
    Sql()
```

**Сырые вырожения в Select**
```go
// SELECT "age", COUNT(*) as count FROM "users" GROUP BY "count"
sql := sqlx.Table("users").
    Select("age", sqlx.Raw("COUNT(*) as count")).
    GroupBy("count").
    Sql()
```

**Вложенный запрос в Select**
```go
// SELECT "name", (SELECT age FROM ag WHERE id = $1 LIMIT 1) as age FROM "users"
sql := sqlx.Table("users").
    Select("name", sqlx.Raw("(SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1)).
    Sql()
```

**Вложенный запрос в Select**
```go
// SELECT "name", (SELECT age FROM ag WHERE id = $1 LIMIT 1) as age FROM "users"
sql := sqlx.Table("users").
    Select("name").
    SelectRaw("(SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1)
    Sql()
```

**Агрегатные функции
Count(column, alias)  
Sum(column, alias)  
Min(column, alias)  
Max(column, alias)  
Avg(column, alias)**
```go
// SELECT "age", COUNT(*) as "count" FROM "users" GROUP BY "count"
sql := sqlx.Table("users").
    Select("age").
    GroupBy("count").
    Count("*", "count")
    Sql()
```

**Условия (Where)**
```go
// SELECT * FROM "users" WHERE "id" = $1
sql := sqlx.Table("users").
    Where("id", "=", 1).
    Sql()
```

**Комбиная условий (And)**
```go
// SELECT * FROM "users" WHERE "age" >= $1 AND "created_at" < $2
sql := sqlx.sqlx.Table("users").
    Where("age", ">=", 21).
    Where("created_at", "<", "2016-01-01").
    Sql().
```

**Комбиная условий (Or)**
```go
// SELECT * FROM "users" WHERE "id" = $1 OR "id" = $2
sql := sqlx.Table("users").
    Where("id", "=", 1).
    OrWhere("id", "=", 2).
    Sql()
```

**Груповые условия Where**
```go
// SELECT * FROM "users" WHERE "group_id" = $1 OR ("age" = $2 OR "age" = $3)
sql := sqlx.Table("users").
    Where("id", "=", 1).
    OrWhereGroup(func(builder *sqlx.Builder) {
        builder.Where("age", "=", 18).OrWhere("age", "=", 21)
    }).Sql()
```

**Сырые условия Where**
```go
// SELECT * FROM "users" WHERE "age" = $1 OR age = $2
sql := sqlx.Table("users").
    Where("age", "=", 18).
    OrWhereRaw("age = ?", 21).
    Sql()
```

**Сырые условия Where с вложенным запросом**
```go
// SELECT * FROM "users" WHERE "age" = $1 OR "age" = (SELECT age FROM ag WHERE id = $2 OR id = $3 LIMIT 1)
sql := sqlx.Table("users").
    Where("age", "=", 18).
    OrWhere("age", "=", sqlx.Raw("(SELECT age FROM ag WHERE id = ? OR id = ? LIMIT 1)", 21, 27)).
    Sql()
```

**Условия (Where Null и Where Not Null)**
```go
// SELECT * FROM "users" WHERE "country" IS NOT NULL
sql := sqlx.Table("users").WhereNotNull("country").Sql()
```

**Условия (Where Between)**
```go
// SELECT * FROM "users" WHERE "create_at" BETWEEN $1 AND $2 AND "create_at" NOT BETWEEN $3 AND $4
sql := sqlx.Table("users").WhereBetween("create_at", "2007-01-01", "2007-12-31").Sql()
```

**Условия (Where In)**
```go
// SELECT * FROM "users" WHERE "id" IN ($1, $2, $3, $4, $5, $6)
sql := sqlx.Table("users").WhereIn("id", sqlx.List{1, 2, 3, 4, 5, 6}).Sql()
```

**Условия Where In c вложенным запросом**
```go
// SELECT * FROM "users" WHERE "id" IN (SELECT "user_id" FROM "orders" WHERE "city" = $1)
sql := sqlx.Table("users").
    WhereIn("id", func(builder *sqlx.Builder) {
        builder.Select("user_id").From("orders").Where("city", "=", "Moscow")
    }).Sql()
```

**Группирока (Group By)**
```go
// SELECT "country", "city", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country", "city"
sql := sqlx.Table("users").
    Select("country", "city").
    Where("age", ">", 18).
    GroupBy("country", "city").
    Count("*", "count").
    Sql()
```

**Сырые группироки**
```go
// SELECT to_char(created_at, 'YYYY-MM-DD') as date, COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY to_char(crated_at, 'YYYY-MM-DD')
sql := sqlx.Table("users").
    SelectRaw("to_char(created_at, 'YYYY-MM-DD') as date").
    Where("age", ">", 18).
    GroupByRaw("to_char(created_at, 'YYYY-MM-DD')").
    Count("*", "count").
    Sql()
```

**Группировка и условия после группировки (HAVING)**
```go
// SELECT "country", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country" HAVING "count" > $2
sql := sqlx.Table("users").
    Select("country").
    Where("age", ">", 18).
    GroupBy("country").
    Having("count", ">", 100).
    Count("*", "count").
    Sql()
```

**Сортировка (ORDER BY)**
```go
// SELECT * FROM "users" WHERE "city" = $1 ORDER BY "id" DESC
sql := sqlx.Table("users").
    Where("city", "=", "Moscow").
    OrderBy("id", "DESC").
    Sql()
```

**Лимиты (Limit и Offset)**
```go
// SELECT * FROM "users" WHERE "city" = $1 ORDER BY "age" ASC LIMIT $2 OFFSET $3
sql := sqlx.Table("users").
    Where("city", "=", "Moscow").
    OrderBy("age", "ASC").
    Limit(10).Offset(50).
    Sql()
```

**Объединение (Join)**
```go
// SELECT * FROM "users" as "us" INNER JOIN "info" as "inf" ON ("us"."id" = "inf"."user_id") WHERE "inf"."city" = $1
sql := sqlx.Table("users as us").
    Join("info as inf", "us.id", "=", "inf.user_id").
    Where("inf.city", "=", "Moscow").
    Sql()
```

**Объединение (Left Join)**
```go
// SELECT * FROM "users" as "us" LEFT JOIN "orders" as "ord" ON ("us"."id" = "ord"."user_id") WHERE "ord"."user_id" IS NOT NULL
sql := sqlx.Table("users as us").
    LeftJoin("orders as ord", "us.id", "=", "ord.user_id").
    WhereNotNull("ord.user_id").
    Sql()
```

**Удаление (Delete)**
```go
// DELETE FROM "users" WHERE "id" = $1
sql := sqlx.Table("users").
    Where("id", "=", 15).
    Delete().
    Sql()
```

**Изменение (Update)**
```go
// UPDATE "users" SET "city" = $1, "name" = $2 WHERE "id" = $3
sql := sqlx.Table("users").
    Where("id", "=", 15).
    Update(sqlx.Data{"name": "Ivan", "city": "Moscow"}).
    Sql()
```

**Вставка (Insert)**
```go
// INSERT INTO "users" ("id", "name") VALUES ($1, $2)
sql := sqlx.Table("users").
    Insert(sqlx.Data{"id": 1, "name": "Jack"}).
    Sql()
```

Если необходимо вставить несколько записей
```go
// INSERT INTO "users" ("id", "name") VALUES ($1, $2), ($3, $4)
sql := sqlx.Table("users").
    Insert(sqlx.Data{"id": 1, "name": "Jack"}, sqlx.Data{"id": 2, "name": "Mike"}).
    Sql()
```

Или так
```go
// INSERT INTO "users" ("id", "name") VALUES ($1, $2), ($3, $4)
sql := sqlx.Table("users").
    Insert(sqlx.Data{"id": 1, "name": "Jack"}).
    Insert(sqlx.Data{"id": 2, "name": "Mike"}).
    Sql()
```

**Выполнение запросов и сканирование результатов**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type User struct {
    Id    int
    Name  string
}

func main() {
    var err error

    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    sqlx.Driver("mysql")
	dbx := sqlx.DataBase(db)

    var users []User

    query := sqlx.Table("users").OrderBy("id", "asc").Limit(10)
	err = dbx.Query(query).Scan(&users)

    if err != nil {
        log.Fatalln("DB Query:", err)
    }

    for _, user := range users {
		fmt.Printf("%d, %s\n", user.Id, user.Name)
	}

	db.Close()
}
```

**Скан в структуру**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type User struct {
    Id    int
    Name  string
}

func main() {
    var err error

    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    sqlx.Driver("mysql")
	dbx := sqlx.DataBase(db)

    var user User

    query := sqlx.Table("users").Where("id", "=", 2).Limit(11)
	err = dbx.Query(query).Scan(&user)

    if err != nil {
        log.Fatalln("DB Query:", err)
    }

	fmt.Printf("%d, %s\n", user.Id, user.Name)

	db.Close()
}
```

**Скан в переменные**

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

func main() {
    var err error

    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    sqlx.Driver("mysql")
	dbx := sqlx.DataBase(db)

    var id int
    var name string

    query := sqlx.Table("users").Select("id", "name").Where("id", "=", 2).Limit(11)
	err = dbx.Query(query).Scan(&id, &name)

    if err != nil {
        log.Fatalln("DB Query:", err)
    }

	fmt.Printf("%d, %s\n", id, name)

	db.Close()
}
```

**Больший контроль над выборкой**

Сканируем и обрабатываем результат по кускам

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/AlexanderGrom/componenta/sqlx"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type User struct {
    Id    int
    Name  string
}

func main() {
    var err error

    db, err := sql.Open("mysql", "UserName:UserPass@/DBName")

    if err != nil {
        log.Fatalln("DB Connecting:", err)
    }

    defer db.Close()

    sqlx.Driver("mysql")
	dbx := sqlx.DataBase(db)

    query := sqlx.Table("users").OrderBy("id", "asc")

    err = dbx.Query(query).Chunk(100, func(users []User) {
		for _, user := range users {
			fmt.Printf("%d, %s\n", user.Id, user.Name)
		}
	})

    if err := rows.Err(); err != nil {
		return log.Fatalln("DB Rows:", err)
	}
}
```
