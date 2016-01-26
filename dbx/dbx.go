package dbx

import (
    "database/sql"
    "fmt"
    "reflect"
)

var (
    ErrNoRows = errors.New("dbx: no rows in result set")
)

type Rows struct {
    *sql.Rows
}

func NewRows(r *sql.Rows) *Rows {
    return &Rows{r}
}

//
// Сканирование в переменные (псевданим Scan)
//
func (self *Rows) ScanVars(a ...interface{}) error {
    return self.Scan(a...)
}

//
// Сканирование в структуру
//
func (self *Rows) ScanStruct(a interface{}) error {
    structValuePrt := reflect.ValueOf(a)

    if structValuePrt.Kind() != reflect.Ptr {
        return errors.New("dbx: destination not a pointe")
    }
    if structValuePrt.IsNil() {
        return errors.New("dbx: destination pointer is nil")
    }

    structValue := reflect.Indirect(structValuePrt)
    structType := structValue.Type()

    if structType.Kind() != reflect.Struct {
        return errors.New("dbx: invalid variable type, must be a struct")
    }

    fields := structType.NumField()
    values := make([]interface{}, fields)

    for i := 0; i < fields; i++ {
        f := structType.Field(i)
        if f.Anonymous || f.PkgPath != "" {
            continue
        }
        values[i] = structValue.Field(i).Addr().Interface()
    }

    return self.Scan(values...)
}

type Fetch struct {
    rows *sql.Rows
}

func NewFetch(r *sql.Rows) *Fetch {
    return &Fetch{r}
}

//
// Сканирование в срез структур
//
func (self *Fetch) ScanSlice(a interface{}) error {
    defer self.rows.Close()

    sliceValuePrt := reflect.ValueOf(a)

    if sliceValuePrt.Kind() != reflect.Ptr {
        return errors.New("dbx: destination not a pointe")
    }
    if sliceValuePrt.IsNil() {
        return errors.New("dbx: destination pointer is nil")
    }

    sliceValue := reflect.Indirect(sliceValuePrt)

    if sliceValue.Kind() != reflect.Slice {
        return errors.New("dbx: invalid variable type, must be a slice")
    }

    for self.rows.Next() {
        structValuePrt := reflect.New(sliceValue.Type().Elem())
        structValue := reflect.Indirect(structValuePrt)

        if structValue.Kind() != reflect.Struct {
            return errors.New("dbx: invalid variable type, must be a slice struct")
        }

        structType := structValue.Type()

        fields := structType.NumField()
        values := make([]interface{}, fields)

        for i := 0; i < fields; i++ {
            f := structType.Field(i)
            if f.Anonymous || f.PkgPath != "" {
                continue
            }
            values[i] = structValue.Field(i).Addr().Interface()
        }

        if err := self.rows.Scan(values...); err != nil {
            return fmt.Errorf("dbx: %s", err)
        }

        sliceValue.Set(reflect.Append(sliceValue, structValue))
    }

    if err := self.rows.Err(); err != nil {
        return fmt.Errorf("dbx: %s", err)
    }

    if sliceValue.Len() == 0 {
        return ErrNoRows
    }

    if err := self.rows.Close(); err != nil {
        return err
    }

    return nil
}

//
// Сканирование в структуру
//
func (self *Fetch) ScanStruct(a interface{}) error {
    defer self.rows.Close()

    structValuePrt := reflect.ValueOf(a)

    if structValuePrt.Kind() != reflect.Ptr {
        return errors.New("dbx: destination not a pointe")
    }
    if structValuePrt.IsNil() {
        return errors.New("dbx: destination pointer is nil")
    }

    structValue := reflect.Indirect(structValuePrt)
    structType := structValue.Type()

    if structType.Kind() != reflect.Struct {
        return errors.New("dbx: invalid variable type, must be a struct")
    }

    if !self.rows.Next() {
        return ErrNoRows
    }

    fields := structType.NumField()
    values := make([]interface{}, fields)

    for i := 0; i < fields; i++ {
        f := structType.Field(i)
        if f.Anonymous || f.PkgPath != "" {
            continue
        }
        values[i] = structValue.Field(i).Addr().Interface()
    }

    if err := self.rows.Scan(values...); err != nil {
        return fmt.Errorf("dbx: %s", err)
    }

    if err := self.rows.Close(); err != nil {
        return err
    }

    return nil
}

//
// Сканирование в переменные (псевдоним Scan)
//
func (self *Fetch) ScanVars(a ...interface{}) error {
    defer self.rows.Close()

    if !self.rows.Next() {
        if err := self.rows.Err(); err != nil {
            return fmt.Errorf("dbx: %s", err)
        }
        return ErrNoRows
    }

    if err := self.rows.Scan(a...); err != nil {
        return fmt.Errorf("dbx: %s", err)
    }

    if err := self.rows.Close(); err != nil {
        return err
    }

    return nil
}
