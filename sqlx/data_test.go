package sqlx

import (
	"testing"
)

// Устанавливаем драйвер для всех тестов
func TestDataDriver(t *testing.T) {
	Driver("pgsql")
}

func TestDataTable(t *testing.T) {
	expect := []interface{}{}
	result := Table("users").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataTable.\nResult: %v\nExpect: %v", result, expect)
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

func TestDataWhereRaw(t *testing.T) {
	expect := []interface{}{}
	result := Table("users").WhereRaw("id = 1").OrWhere("age", "=", Raw("2")).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereRaw.\nResult: %v\nExpect: %v", result, expect)
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

func TestDataWhereInSub(t *testing.T) {
	expect := []interface{}{"Moscow"}
	result := Table("users").WhereIn("id", func(builder *Builder) {
		builder.Select("user_id").From("orders").Where("city", "=", "Moscow")
	}).Select("*").Data()
	if !DataEqual(result, expect) {
		t.Errorf("Expect result to equal in func TestDataWhereInSub.\nResult: %v\nExpect: %v", result, expect)
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
