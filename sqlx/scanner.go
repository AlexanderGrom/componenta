package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

var ErrNoRows = errors.New("sqlx: no rows in result set")

type Scanner struct {
	rows *sql.Rows
}

func NewScanner(r *sql.Rows) *Scanner {
	return &Scanner{r}
}

// Сканирование результатов запроса в переменные
func (self *Scanner) Scan(a ...interface{}) error {
	defer self.rows.Close()

	if len(a) == 0 {
		return errors.New("sqlx: no destination")
	}

	switch reflect.Indirect(reflect.ValueOf(a[0])).Type().Kind() {
	case reflect.Struct:
		return self.scanStruct(a[0])
	case reflect.Slice:
		return self.scanSlice(a[0])
	default:
		return self.scanVars(a...)
	}
}

// Сканирование в срез структур
func (self *Scanner) scanSlice(a interface{}) error {
	sliceValuePrt := reflect.ValueOf(a)

	if sliceValuePrt.Kind() != reflect.Ptr {
		return errors.New("sqlx: destination not a pointe")
	}
	if sliceValuePrt.IsNil() {
		return errors.New("sqlx: destination pointer is nil")
	}

	sliceValue := reflect.Indirect(sliceValuePrt)

	if sliceValue.Kind() != reflect.Slice {
		return errors.New("sqlx: invalid variable type, must be a slice")
	}

	columns, err := self.rows.Columns()

	if err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	sliceElem := sliceValue.Type().Elem()
	structValuePrt := reflect.New(sliceElem)
	structValue := reflect.Indirect(structValuePrt)

	if structValue.Kind() != reflect.Struct {
		return errors.New("sqlx: invalid variable type, must be a slice struct")
	}

	structType := structValue.Type()

	count := structType.NumField()
	fieldsMap := make(map[int]string)

	for i := 0; i < count; i++ {
		f := structType.Field(i)
		if f.Anonymous || f.PkgPath != "" {
			continue
		}
		fieldsMap[i] = toSnake(f.Name)
	}

	for self.rows.Next() {
		structValuePrt := reflect.New(sliceElem)
		structValue := reflect.Indirect(structValuePrt)

		fields := make(map[string]reflect.Value)

		for k, v := range fieldsMap {
			fields[v] = structValue.Field(k)
		}

		values := make([]interface{}, len(columns))

		for i, column := range columns {
			if field, ok := fields[column]; ok {
				values[i] = field.Addr().Interface()
			} else {
				values[i] = &sql.NullString{}
			}
		}

		if err := self.rows.Scan(values...); err != nil {
			return fmt.Errorf("sqlx: %s", err)
		}

		sliceValue.Set(reflect.Append(sliceValue, structValue))
	}

	if err := self.rows.Err(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	if sliceValue.Len() == 0 {
		return ErrNoRows
	}

	if err := self.rows.Close(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	return nil
}

// Сканирование в структуру
func (self *Scanner) scanStruct(a interface{}) error {
	structValuePrt := reflect.ValueOf(a)

	if structValuePrt.Kind() != reflect.Ptr {
		return errors.New("sqlx: destination not a pointe")
	}
	if structValuePrt.IsNil() {
		return errors.New("sqlx: destination pointer is nil")
	}

	structValue := reflect.Indirect(structValuePrt)
	structType := structValue.Type()

	if structType.Kind() != reflect.Struct {
		return errors.New("sqlx: invalid variable type, must be a struct")
	}

	if !self.rows.Next() {
		if err := self.rows.Err(); err != nil {
			return fmt.Errorf("sqlx: %s", err)
		}
		return ErrNoRows
	}

	columns, err := self.rows.Columns()

	if err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	count := structType.NumField()
	fields := make(map[string]reflect.Value)

	for i := 0; i < count; i++ {
		f := structType.Field(i)
		if f.Anonymous || f.PkgPath != "" {
			continue
		}
		fields[toSnake(f.Name)] = structValue.Field(i)
	}

	values := make([]interface{}, len(columns))

	for i, column := range columns {
		if field, ok := fields[column]; ok {
			values[i] = field.Addr().Interface()
		} else {
			values[i] = &sql.NullString{}
		}
	}

	if err := self.rows.Scan(values...); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	if err := self.rows.Close(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	return nil
}

// Сканирование в переменные
func (self *Scanner) scanVars(a ...interface{}) error {
	if !self.rows.Next() {
		if err := self.rows.Err(); err != nil {
			return fmt.Errorf("sqlx: %s", err)
		}
		return ErrNoRows
	}

	if err := self.rows.Scan(a...); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	if err := self.rows.Close(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	return nil
}
