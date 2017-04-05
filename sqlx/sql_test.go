package sqlx

import (
	"testing"
)

// Устанавливаем драйвер для всех тестов
func TestSqlDriver(t *testing.T) {
	Driver("postgres")
}

func TestSqlTable1(t *testing.T) {
	expect := `SELECT * FROM "users"`
	result := Table("users").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlTable1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlTable2(t *testing.T) {
	expect := `SELECT * FROM ( SELECT "group_id", MAX("created_at") as "lastdate" FROM "users" GROUP BY "group_id" ) as "users" ORDER BY "lastdate" DESC`
	result := Table(func(builder *Builder) {
		builder.Select("group_id").From("users").GroupBy("group_id").Max("created_at", "lastdate")
	}).OrderBy("lastdate", "DESC").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlTable2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlTable3(t *testing.T) {
	expect := `SELECT * FROM ( SELECT * FROM "users" WHERE "id" = $1 ) as users`
	subque := Table("users").Where("id", "=", 1)
	result := Table(Raw("( "+subque.Sql()+" ) as users", subque.Data()...)).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlTable3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlTable4(t *testing.T) {
	expect := `SELECT * FROM ( SELECT * FROM "users" WHERE "id" = $1 ) as "users"`
	subque := Table("users").Where("id", "=", 1)
	result := Table(subque).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlTable4.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlTable5(t *testing.T) {
	expect := `SELECT * FROM (SELECT * FROM users WHERE id IN ($1, $2, $3)) as users WHERE "age" > $4`
	result := Table(Raw("(SELECT * FROM users WHERE id IN (?, ?, ?)) as users", 1, 2, 3)).Where("age", ">", 21).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlTable5.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelect1(t *testing.T) {
	expect := `SELECT "id" FROM "users"`
	result := Table("users").Select("id").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelect1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelect2(t *testing.T) {
	expect := `SELECT "id", "name", "age" FROM "users"`
	result := Table("users").Select("id", "name", "age").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelect2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelect3(t *testing.T) {
	expect := `SELECT "id", "name", "age", COUNT(*) as count FROM "users"`
	result := Table("users").Select("id", "name", "age", Raw("COUNT(*) as count")).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelect3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelect4(t *testing.T) {
	expect := `SELECT "id", "name", "age", COUNT(*) as count FROM "users"`
	result := Table("users").Select("id", "name").Select("age").SelectRaw("COUNT(*) as count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelect4.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelectRaw1(t *testing.T) {
	expect := `SELECT COUNT(*) as count FROM "users"`
	result := Table("users").Select(Raw("COUNT(*) as count")).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelectRaw1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelectRaw2(t *testing.T) {
	expect := `SELECT COUNT(*) as count FROM "users"`
	result := Table("users").SelectRaw("COUNT(*) as count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelectRaw2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelectRaw3(t *testing.T) {
	expect := `SELECT "name", (SELECT age FROM ag WHERE id = $1 LIMIT 1) as age FROM "users"`
	result := Table("users").Select("name", Raw("(SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1)).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelectRaw3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlSelectRaw4(t *testing.T) {
	expect := `SELECT name, (SELECT age FROM ag WHERE id = $1 LIMIT 1) as age FROM "users"`
	result := Table("users").SelectRaw("name, (SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlSelectRaw4.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereBase1(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" = $1`
	result := Table("users").Where("id", "=", 1).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereBase1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereBase2(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" = $1 OR "id" = $2`
	result := Table("users").Where("id", "=", 1).OrWhere("id", "=", 2).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereBase2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereBase3(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" = $1 AND "age" = $2`
	result := Table("users").Where("id", "=", 1).Where("age", "=", 31).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereBase3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereGroup(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" = $1 OR ( "id" = $2 OR "id" = $3 )`
	result := Table("users").Where("id", "=", 1).OrWhereGroup(func(builder *Builder) {
		builder.Where("id", "=", 2).OrWhere("id", "=", 3)
	}).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereGroup.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereRaw1(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE id = 1 OR "age" = 2`
	result := Table("users").WhereRaw("id = 1").OrWhere("age", "=", Raw("2")).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereRaw1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereRaw2(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "age" = $1 OR age = $2`
	result := Table("users").Where("age", "=", 1).OrWhereRaw("age = ?", 2).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereRaw2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereRaw3(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "age" = $1 OR "age" = (SELECT age FROM ag WHERE id = $2 OR id = $3 LIMIT 1)`
	result := Table("users").Where("age", "=", 1).OrWhere("age", "=", Raw("(SELECT age FROM ag WHERE id = ? OR id = ? LIMIT 1)", 2, 3)).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereRaw3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereNull(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "country" IS NULL AND "city" IS NOT NULL`
	result := Table("users").WhereNull("country").WhereNotNull("city").Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereNull.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereBETWEEN(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "create_at" BETWEEN $1 AND $2 AND "create_at" NOT BETWEEN $3 AND $4`
	result := Table("users").
		WhereBetween("create_at", "2007-01-01", "2007-12-31").
		WhereNotBetween("create_at", "2007-02-01", "2007-02-31").
		Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereBETWEEN.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereInList(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" IN ( $1, $2, $3, $4, $5, $6 )`
	result := Table("users").WhereIn("id", List{1, 2, 3, 4, 5, 6}).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereInList.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereInSub1(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "id" IN ( SELECT "user_id" FROM "orders" WHERE "city" = $1 )`
	result := Table("users").WhereIn("id", func(builder *Builder) {
		builder.Select("user_id").From("orders").Where("city", "=", "Moscow")
	}).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereInSub1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereInSub2(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "deleted" = $1 AND "id" IN ( SELECT "user_id" FROM "orders" WHERE "city" = $2 )`
	result := Table("users").Where("deleted", "=", "1").
		WhereIn("id", func(builder *Builder) {
			builder.Select("user_id").From("orders").Where("city", "=", "Moscow")
		}).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereInSub2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlWhereInSub3(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "deleted" = $1 AND "id" IN ( SELECT "user_id" FROM "orders" WHERE "city" = $2 )`
	subque := Table("orders").Where("city", "=", "Moscow").Select("user_id")
	result := Table("users").Where("deleted", "=", "1").WhereIn("id", subque).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlWhereInSub3.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlGroupBy(t *testing.T) {
	expect := `SELECT "country", "city", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country", "city"`
	result := Table("users").Where("age", ">", 31).GroupBy("country", "city").Select("country", "city").Count("*", "count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlGroupBy.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlHavingBase(t *testing.T) {
	expect := `SELECT "country", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country" HAVING "count" > $2`
	result := Table("users").Where("age", ">", 31).GroupBy("country").Having("count", ">", 10).Select("country").Count("*", "count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlHavingBase.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlHavingRaw(t *testing.T) {
	expect := `SELECT "country", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country" HAVING count > 10`
	result := Table("users").Where("age", ">", 31).GroupBy("country").HavingRaw("count > 10").Select("country").Count("*", "count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlHavingRaw.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlHavingGroup(t *testing.T) {
	expect := `SELECT "country", COUNT(*) as "count" FROM "users" WHERE "age" > $1 GROUP BY "country" HAVING ( "count" > $2 AND "count" < $3 )`
	result := Table("users").Where("age", ">", 31).GroupBy("country").HavingGroup(func(builder *Builder) {
		builder.Having("count", ">", 10).Having("count", "<", 20)
	}).Select("country").Count("*", "count").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlHavingGroup.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlOrderBy1(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "city" = $1 ORDER BY "id" DESC`
	result := Table("users").Where("city", "=", "Moscow").OrderBy("id", "DESC").Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlOrderBy1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlOrderBy2(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "city" = $1 ORDER BY "age" ASC, "name" DESC`
	result := Table("users").Where("city", "=", "Moscow").OrderBy("age", "ASC").OrderBy("name", "DESC").Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlOrderBy2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlLimit(t *testing.T) {
	expect := `SELECT * FROM "users" WHERE "city" = $1 ORDER BY "age" ASC LIMIT $2 OFFSET $3`
	result := Table("users").Where("city", "=", "Moscow").OrderBy("age", "ASC").Limit(10).Offset(50).Select("*").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlLimit.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlJoin(t *testing.T) {
	expect := `SELECT * FROM "users" as "us" INNER JOIN "info" as "inf" ON ( "us"."id" = "inf"."user_id" AND "us"."group" = $1 ) INNER JOIN "city" as "c" ON ( "c"."id" = "inf"."city_id" )`
	result := Table("users as us").Join("info as inf", func(joiner *Joiner) {
		joiner.On("us.id", "=", "inf.user_id")
		joiner.Where("us.group", "=", "admin")
	}).Join("city as c", func(joiner *Joiner) {
		joiner.On("c.id", "=", "inf.city_id")
	}).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlJoin.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlLeftJoin(t *testing.T) {
	expect := `SELECT * FROM "users" as "us" LEFT JOIN "orders" as "ord" ON ( "us"."id" = "ord"."user_id" ) WHERE "ord"."user_id" IS NOT NULL`
	result := Table("users as us").LeftJoin("orders as ord", func(joiner *Joiner) {
		joiner.On("us.id", "=", "ord.user_id")
	}).WhereNotNull("ord.user_id").Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlLeftJoin.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlDelete(t *testing.T) {
	expect := `DELETE FROM "users" WHERE "id" = $1`
	result := Table("users").Where("id", "=", 15).Delete().Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlDelete.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlUpdate(t *testing.T) {
	expect := `UPDATE "users" SET "city" = $1, "name" = $2 WHERE "id" = $3`
	result := Table("users").Where("id", "=", 15).Update(Data{"name": "Jack", "city": "Moscow"}).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlUpdate.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlInsert1(t *testing.T) {
	expect := `INSERT INTO "users" ( "id", "name" ) VALUES ( $1, $2 )`
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlInsert1.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlInsert2(t *testing.T) {
	expect := `INSERT INTO "users" ( "id", "name" ) VALUES ( $1, $2 ), ( $3, $4 )`
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}, Data{"id": 2, "name": "Mike"}).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlInsert2.\nResult: %s\nExpect: %s", result, expect)
	}
}

func TestSqlInsert3(t *testing.T) {
	expect := `INSERT INTO "users" ( "id", "name" ) VALUES ( $1, $2 ), ( $3, $4 )`
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}).Insert(Data{"id": 2, "name": "Mike"}).Sql()
	if result != expect {
		t.Errorf("Expect result to equal in func TestSqlInsert3.\nResult: %s\nExpect: %s", result, expect)
	}
}
