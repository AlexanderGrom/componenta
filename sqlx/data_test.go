package sqlx

import (
	"testing"
)

// Устанавливаем драйвер для всех тестов
func TestDataDriver(t *testing.T) {
	Driver("postgres")
}

func TestDataTable1(t *testing.T) {
	expect := []interface{}{}
	result := Table("users").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataTable1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataTable2(t *testing.T) {
	expect := []interface{}{1, 2, 3, 21}
	result := Table(func(builder *Builder) {
		builder.Select("*").From("users").WhereIn("id", List{1, 2, 3})
	}).Where("age", ">", 21).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataTable2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataTable3(t *testing.T) {
	expect := []interface{}{1, 2, 3, 21}
	subque := Table("users").WhereIn("id", List{1, 2, 3})
	result := Table(subque).Where("age", ">", 21).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataTable3.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataTable4(t *testing.T) {
	expect := []interface{}{1, 2, 3, 21}
	result := Table(Raw("(SELECT * FROM users WHERE id IN (?, ?, ?) as users", 1, 2, 3)).Where("age", ">", 21).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataTable4.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataSelect1(t *testing.T) {
	expect := []interface{}{1, 2}
	result := Table("users").Select("name", Raw("(SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1)).Where("id", "=", 2).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataSelect1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataSelect2(t *testing.T) {
	expect := []interface{}{1, 2}
	result := Table("users").SelectRaw("(SELECT age FROM ag WHERE id = ? LIMIT 1) as age", 1).Where("id", "=", 2).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataSelect2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereBase1(t *testing.T) {
	expect := []interface{}{1}
	result := Table("users").Where("id", "=", 1).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereBase1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereBase2(t *testing.T) {
	expect := []interface{}{1, 2}
	result := Table("users").Where("id", "=", 1).OrWhere("id", "=", 2).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereBase2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereGroup(t *testing.T) {
	expect := []interface{}{1, 2, 3}
	result := Table("users").Where("id", "=", 1).OrWhereGroup(func(builder *Builder) {
		builder.Where("id", "=", 2).OrWhere("id", "=", 3)
	}).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereGroup.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereRaw1(t *testing.T) {
	expect := []interface{}{}
	result := Table("users").WhereRaw("id = 1").OrWhere("age", "=", Raw("2")).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereRaw1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereRaw2(t *testing.T) {
	expect := []interface{}{1, 2, 3}
	result := Table("users").Where("age", "=", 1).OrWhereRaw("age = ?", 2).OrWhere("age", "=", Raw("(SELECT age FROM users WHERE id = ? LIMIT 1)", 3)).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereRaw2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereBETWEEN(t *testing.T) {
	expect := []interface{}{"2007-01-01", "2007-12-31", "2007-02-01", "2007-02-31"}
	result := Table("users").
		WhereBetween("create_at", "2007-01-01", "2007-12-31").
		WhereNotBetween("create_at", "2007-02-01", "2007-02-31").
		Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereBETWEEN.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereInList(t *testing.T) {
	expect := []interface{}{1, 2, 3, 4, 5, 6}
	result := Table("users").WhereIn("id", List{1, 2, 3, 4, 5, 6}).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereInList.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereInSub1(t *testing.T) {
	expect := []interface{}{"Moscow"}
	result := Table("users").WhereIn("id", func(builder *Builder) {
		builder.Select("user_id").From("orders").Where("city", "=", "Moscow")
	}).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereInSub1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataWhereInSub2(t *testing.T) {
	expect := []interface{}{"Moscow"}
	subque := Table("orders").Where("city", "=", "Moscow").Select("user_id")
	result := Table("users").WhereIn("id", subque).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereInSub2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataHavingBase(t *testing.T) {
	expect := []interface{}{31, 10}
	result := Table("users").Where("age", ">", 31).GroupBy("country").Having("count", ">", 10).Select("country").Count("*", "count").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataHavingBase.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataHavingRaw(t *testing.T) {
	expect := []interface{}{31}
	result := Table("users").Where("age", ">", 31).GroupBy("country").HavingRaw("count > 10").Select("country").Count("*", "count").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataHavingRaw.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataHavingGroup(t *testing.T) {
	expect := []interface{}{31, 10, 20}
	result := Table("users").Where("age", ">", 31).GroupBy("country").HavingGroup(func(builder *Builder) {
		builder.Having("count", ">", 10).Having("count", "<", 20)
	}).Select("country").Count("*", "count").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataHavingGroup.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataLimit(t *testing.T) {
	expect := []interface{}{"Moscow", 10, 50}
	result := Table("users").Where("city", "=", "Moscow").OrderBy("age", "ASC").Limit(10).Offset(50).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataLimit.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataJoin(t *testing.T) {
	expect := []interface{}{"admin", 1000, 10}
	result := Table("users as us").Join("info as inf", func(joiner *Joiner) {
		joiner.On("us.id", "=", "inf.user_id")
		joiner.Where("us.group", "=", "admin")
	}).Where("us.id", ">", 1000).Limit(10).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataJoin.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataUpdate(t *testing.T) {
	expect := []interface{}{"Moscow", "Jack", 15}
	result := Table("users").Where("id", "=", 15).Update(Data{"name": "Jack", "city": "Moscow"}).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataUpdate.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataInsert1(t *testing.T) {
	expect := []interface{}{1, "Jack"}
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataInsert1.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataInsert2(t *testing.T) {
	expect := []interface{}{1, "Jack", 2, "Mike"}
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}, Data{"id": 2, "name": "Mike"}).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataInsert2.\nResult: %v\nExpect: %v", result, expect)
	}
}

func TestDataInsert3(t *testing.T) {
	expect := []interface{}{1, "Jack", 2, "Mike"}
	result := Table("users").Insert(Data{"id": 1, "name": "Jack"}).Insert(Data{"id": 2, "name": "Mike"}).Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataInsert3.\nResult: %v\nExpect: %v", result, expect)
	}
}

func DataEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
