package sqlx

import (
	"fmt"
	"strconv"
)

// Интерерфейс Stringer
type Stringer interface {
	String() string
}

// Конвертирует интерфейс в строку
func toString(x interface{}) string {
	switch x := x.(type) {
	case string:
		return x
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(int64(x), 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(uint64(x), 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', 6, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', 6, 64)
	case Stringer:
		return x.String()
	case []byte:
		return string(x)
	case []rune:
		return string(x)
	case bool:
		if x {
			return "1"
		}
		return "0"
	case nil:
		return ""
	}
	// Sprint довольно тяжелый метод, поэтому мы делаем некоторые очевидные
	// проверки сами
	return fmt.Sprint(x)
}

// Комбинирует строки в одну строку
func combine(x ...interface{}) string {
	if len(x) == 0 {
		return ""
	}
	a := make([]string, len(x))
	for k, v := range x {
		a[k] = toString(v)
	}
	x = nil
	if len(a) == 1 {
		return a[0]
	}
	sep := " "
	num := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		num += len(a[i])
	}
	b := make([]byte, num)
	var n int
	var nn int
	for _, s := range a {
		if nn > 0 {
			n += copy(b[n:], sep)
		}
		nn = copy(b[n:], s)
		n += nn
	}
	for n > 0 {
		if b[n-1] != ' ' {
			break
		}
		n--
	}
	return string(b[:n])
}

// Проверяет оператор на корректность
func isOperator(o string) bool {
	switch o {
	case "=", "!=", "<", ">", "<=", ">=", "<>",
		"&", "|", "^", "<<", ">>",
		"LIKE", "NOT LIKE":
		return true
	}
	return false
}
