package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Конфиг
type Config map[string]interface{}

// Вернет значение или nil если знаение не найдено
func (self Config) Get(name string) interface{} {
	x := self[name]
	if x == nil {
		return nil
	}
	switch a := x.(type) {
	case []interface{}:
		if len(a) > 0 {
			return a[0]
		}
	default:
		return x
	}
	return nil
}

// Вернет число
func (self Config) GetInt(name string) int64 {
	x := self.Get(name)
	if x == nil {
		return 0
	}
	a, _ := strconv.ParseInt(fmt.Sprint(x), 10, 64)
	return a
}

// Вернет флоат
func (self Config) GetFloat(name string) float64 {
	x := self.Get(name)
	if x == nil {
		return 0
	}
	a, _ := strconv.ParseFloat(fmt.Sprint(x), 64)
	return a
}

// Вернет строку
func (self Config) GetString(name string) string {
	x := self.Get(name)
	if x == nil {
		return ""
	}
	return fmt.Sprint(x)
}

// Вернет буул
func (self Config) GetBool(name string) bool {
	x := self.GetString(name)
	if x == "" {
		return false
	}
	switch strings.ToLower(x) {
	case "true", "1", "on", "yes":
		return true
	case "false", "0", "off", "no":
		return false
	}
	return false
}

// Вернет срез или nil если значение не найдено
func (self Config) GetSlice(name string) []interface{} {
	x := self[name]
	if x == nil {
		return nil
	}
	switch a := x.(type) {
	case []interface{}:
		return a
	default:
		return []interface{}{x}
	}
}

// Вернет срез чисел
func (self Config) GetInts(name string) []int64 {
	x := self.GetSlice(name)
	if x == nil {
		return []int64{}
	}
	s := make([]int64, len(x))
	for i, v := range x {
		a, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		s[i] = a
	}
	return s
}

// Вернет срез флоатов
func (self Config) GetFloats(name string) []float64 {
	x := self.GetSlice(name)
	if x == nil {
		return []float64{}
	}
	s := make([]float64, len(x))
	for i, v := range x {
		a, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
		s[i] = a
	}
	return s
}

// Вернет срез строк
func (self Config) GetStrings(name string) []string {
	x := self.GetSlice(name)
	if x == nil {
		return []string{}
	}
	s := make([]string, len(x))
	for i, v := range x {
		s[i] = fmt.Sprint(v)
	}
	return s
}

// Вернет срез логических значений
func (self Config) GetBools(name string) []bool {
	x := self.GetSlice(name)
	if x == nil {
		return []bool{}
	}
	s := make([]bool, len(x))
	for i, v := range x {
		switch strings.ToLower(v.(string)) {
		case "true", "1", "on", "yes":
			s[i] = true
		case "false", "0", "off", "no":
			s[i] = false
		default:
			s[i] = false
		}
	}
	return s
}

// Проверка на сущестование значения
func (self Config) Exists(name string) bool {
	if _, ok := self[name]; !ok {
		return false
	}
	return true
}
