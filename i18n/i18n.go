package i18n

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

//
// Хранилище
//
var storage map[string]Lang

//
// Ошибки
//
var (
    ErrPathEmpty      = errors.New("i18n: path empty")
    ErrPathNotCorrect = errors.New("i18n: path not correct")
    ErrLangDontOpen   = errors.New("i18n: lang file don`t open")
)

//
// Использование конкретного языкового файла
//
func Use(path string) (Lang, error) {
    if len(path) == 0 {
        return Lang{}, ErrPathEmpty
    }
    path, err := filepath.Abs(path)
    if err != nil {
        return Lang{}, ErrPathNotCorrect
    }
    if lng, ok := storage[path]; ok {
        return lng, nil
    }
    file, err := os.Open(path)
    if err != nil {
        return Lang{}, ErrLangDontOpen
    }
    defer file.Close()
    lang := parse(file)
    return lang, nil
}

//
// Парсинг файла
//
func parse(r io.Reader) Lang {
    lang := Lang{}
    scanner := bufio.NewScanner(r)
    parser := newParser()
    for scanner.Scan() {
        key, value, matched := parser.Parse(scanner.Text())
        if matched {
            lang[key] = value
        }
    }
    return lang
}

//
// Языковая настройка
//
type Lang map[string]interface{}

//
// Вернет значение или nil если знаение не найдено
//
func (self Lang) Get(name string) interface{} {
    x := self[name]
    if x == nil {
        return nil
    }
    switch x.(type) {
    case []interface{}:
        return nil
    default:
        return x
    }
}

//
// Вернет строку
//
func (self Lang) GetString(name string) string {
    x := self.Get(name)
    if x == nil {
        return ""
    }
    return fmt.Sprint(x)
}

//
// Вернет срез интерфейсов
//
func (self Lang) GetSlice(name string) []interface{} {
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

//
// Вернет срез строк
//
func (self Lang) GetStrings(name string) []string {
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

//
// Проверка на существование
//
func (self Lang) Exists(name string) bool {
    if _, ok := self[name]; !ok {
        return false
    }
    return true
}
