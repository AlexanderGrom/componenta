package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type ChunkFunk interface{}

type Chunker struct {
	rows *sql.Rows
}

func NewChunker(r *sql.Rows) *Chunker {
	return &Chunker{r}
}

func (self *Chunker) Chunk(n int, f ChunkFunk) error {
	defer self.rows.Close()

	funcValue := reflect.ValueOf(f)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		return errors.New("sqlx: two parameter not a func")
	}

	if funcType.NumIn() != 1 {
		return errors.New("sqlx: two parameter must be func([]Struct)")
	}

	sliceType := funcType.In(0)

	if sliceType.Kind() != reflect.Slice {
		return errors.New("sqlx: two parameter must be func([]Struct)")
	}

	numOut := funcType.NumOut()

	if numOut > 1 {
		return errors.New("sqlx: return must be bool")
	}

	if numOut == 1 {
		outType := funcType.Out(0)
		if outType.Kind() != reflect.Bool {
			return errors.New("sqlx: return must be bool")
		}
	}

	columns, err := self.rows.Columns()

	if err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	sliceValuePrt := reflect.New(sliceType)
	sliceValue := reflect.Indirect(sliceValuePrt)

	sliceElem := sliceType.Elem()
	structValuePrt := reflect.New(sliceElem)
	structValue := reflect.Indirect(structValuePrt)

	if structValue.Kind() != reflect.Struct {
		return errors.New("sqlx: two parameter must be func([]Struct)")
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

	var found bool = false
	var interrupt bool = false
	var i int = 1
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

		if i == n {
			out := funcValue.Call([]reflect.Value{sliceValue})
			sliceValue = reflect.Indirect(reflect.New(sliceType))
			found = true
			i = 1
			if len(out) > 0 && !out[0].Bool() {
				interrupt = true
				break
			}
			continue
		}

		i++
		found = true
	}

	if !interrupt && sliceValue.Len() > 0 {
		funcValue.Call([]reflect.Value{sliceValue})
	}

	if err := self.rows.Err(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	if !found {
		return sql.ErrNoRows
	}

	if err := self.rows.Close(); err != nil {
		return fmt.Errorf("sqlx: %s", err)
	}

	return nil
}
