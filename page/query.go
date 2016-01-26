package page

import (
    "bytes"
    "net/url"
    "strconv"
    "strings"
)

func (self *Page) Query(name string) *Result {
    pathString := "/"
    queryString := ""
    if uri, err := url.ParseRequestURI(self.result.CurrentURI); err == nil {
        pathString = uri.EscapedPath()
        queryString = uri.RawQuery
    }

    query, _ := parseQuery(queryString)

    for k, v := range self.result.List {
        if v["num"] == "1" {
            query := query.Copy()
            query.Del(name)
            if queryString != "" {
                self.result.List[k]["url"] = pathString + "?" + query.Encode()
            } else {
                self.result.List[k]["url"] = pathString
            }
        } else {
            query := query.Copy()
            query.Set(name, v["num"])
            self.result.List[k]["url"] = pathString + "?" + query.Encode()
        }
    }

    if self.result.Current == 1 || self.result.Prev == 1 {
        query := query.Copy()
        query.Del(name)
        if queryString != "" {
            self.result.PrevURI = pathString + "?" + query.Encode()
        } else {
            self.result.PrevURI = pathString
        }
    } else {
        query := query.Copy()
        query.Set(name, strconv.Itoa(self.result.Prev))
        self.result.PrevURI = pathString + "?" + query.Encode()
    }

    if self.result.Current == 1 {
        query := query.Copy()
        query.Del(name)
        if queryString != "" {
            self.result.CurrentURI = pathString + "?" + query.Encode()
        } else {
            self.result.CurrentURI = pathString
        }
    } else {
        query := query.Copy()
        query.Set(name, strconv.Itoa(self.result.Current))
        self.result.CurrentURI = pathString + "?" + query.Encode()
    }

    if self.result.Current == self.result.Total {
        query := query.Copy()
        query.Set(name, strconv.Itoa(self.result.Current))
        self.result.NextURI = pathString + "?" + query.Encode()
    } else {
        query := query.Copy()
        query.Set(name, strconv.Itoa(self.result.Next))
        self.result.NextURI = pathString + "?" + query.Encode()
    }

    if queryString != "" {
        query := query.Copy()
        query.Del(name)
        self.result.FirstURI = pathString + "?" + query.Encode()
    } else {
        self.result.FirstURI = pathString
    }

    if self.result.Total > 1 {
        query := query.Copy()
        query.Set(name, strconv.Itoa(self.result.Total))
        self.result.LastURI = pathString + "?" + query.Encode()
    } else {
        self.result.LastURI = self.result.FirstURI
    }

    return self.result
}

//
// Переписываем стандартный url.ParseQuery() и тип url.Value
// Новый тип сохраняет порядок следования параметров при разборе, изменении и обратной сборке
// Порядок следования параметров должен быть такой же, который пришел от клиента
//
type query []map[string]string

//
// Парсит строку запроса
//
func parseQuery(s string) (query, error) {
    q := make(query, 0)
    for s != "" {
        key := s
        if i := strings.IndexAny(key, "&;"); i >= 0 {
            key, s = key[:i], key[i+1:]
        } else {
            s = ""
        }
        if key == "" {
            continue
        }
        value := ""
        if i := strings.Index(key, "="); i >= 0 {
            key, value = key[:i], key[i+1:]
        }
        key, err := url.QueryUnescape(key)
        if err != nil {
            return q, err
        }
        value, err = url.QueryUnescape(value)
        if err != nil {
            return q, err
        }
        q = append(q, map[string]string{
            key: value,
        })
    }
    return q, nil
}

//
// Возвращает параметр
//
func (self *query) Get(key string) string {
    if *self == nil {
        return ""
    }
    for _, m := range *self {
        if v, ok := m[key]; ok {
            return v
        }
    }
    return ""
}

//
// Изменяет параметр или добавляет новый если такого нет
//
func (self *query) Set(key, value string) {
    found := false
    for _, m := range *self {
        if _, ok := m[key]; ok {
            m[key] = value
            found = true
            break
        }
    }
    if !found {
        self.Add(key, value)
    }
}

//
// Добавляет новый параметр
//
func (self *query) Add(key, value string) {
    *self = append(*self, map[string]string{
        key: value,
    })
}

//
// Удаляет параметр
//
func (self *query) Del(key string) {
    for k, m := range *self {
        if _, ok := m[key]; ok {
            *self = append((*self)[:k], (*self)[k+1:]...)
        }
    }
}

//
// Кодирует query и возвращает строку
//
func (self *query) Encode() string {
    var buf bytes.Buffer
    for _, m := range *self {
        for k, v := range m {
            if buf.Len() > 0 {
                buf.WriteByte('&')
            }
            buf.WriteString(url.QueryEscape(k))
            buf.WriteString("=")
            buf.WriteString(url.QueryEscape(v))
        }
    }
    return buf.String()
}

//
// Возвращает копию
//
func (self *query) Copy() query {
    return append(query(nil), *self...)
}
