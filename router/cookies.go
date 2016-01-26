package router

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "net/http"
    "strings"
)

// Серкетный ключ для присоединения к значению кука хеша
// С целью проверки на изменение значения на клиенте
var (
    TOKEN = "wed5623015c29280c944e0dd937f8995"
)

//
// CookieWriter получает http.ResponseWriter и служит созданием куков для ответа
//
type CookieWriter struct {
    w http.ResponseWriter
}

func NewCookieWriter(w http.ResponseWriter) *CookieWriter {
    return &CookieWriter{w}
}

//
// Установка куки
//
func (self *CookieWriter) Set(key, value string, age int) {
    http.SetCookie(self.w, &http.Cookie{
        Name:     key,
        Value:    self.hash(value) + "+" + value,
        Path:     "/",
        MaxAge:   age,
        HttpOnly: true,
    })
}

//
// Установка сырой куки (без приписанного хеша)
//
func (self *CookieWriter) SetRaw(cookie *http.Cookie) {
    http.SetCookie(self.w, cookie)
}

//
// Удаление куки
//
func (self *CookieWriter) Del(key string) {
    http.SetCookie(self.w, &http.Cookie{
        Name:   key,
        Path:   "/",
        MaxAge: -1,
    })
}

//
// Создание хеша, который будет приписан к значению куки
// Целью хеша является ключ для проверки приходящих куков на предмет их модификации клиентом
//
func (self *CookieWriter) hash(value string) string {
    h := hmac.New(sha256.New, []byte(TOKEN))
    h.Write([]byte(value))
    return fmt.Sprintf("%x", h.Sum(nil))
}

//
// CookieReader принимает http.Request, и служит для чтению куков пришедших от клиента
//
type CookieReader struct {
    data map[string]*http.Cookie
}

func NewCookieReader(r *http.Request) *CookieReader {
    c := &CookieReader{
        data: make(map[string]*http.Cookie),
    }
    for _, v := range r.Cookies() {
        c.data[v.Name] = v
    }
    return c
}

//
// Получения значения куки, установленной через ctx.Res.Cookies.Set
//
func (self *CookieReader) Get(key string) string {
    if self.Exists(key) {
        return self.parse(self.data[key].Value)
    }
    return ""
}

//
// Проверка куки на существование
//
func (self *CookieReader) Exists(key string) bool {
    _, ok := self.data[key]
    return ok
}

//
// Получение сырого значения куки, нужно если кука устанавливается на клиенте
// и не имеет приписаного хеша валадации
//
func (self *CookieReader) GetRaw(key string) *http.Cookie {
    return self.data[key]
}

//
// Создание хеша, для проверки куки на предмет модицикации
//
func (self *CookieReader) hash(value string) string {
    h := hmac.New(sha256.New, []byte(TOKEN))
    h.Write([]byte(value))
    return fmt.Sprintf("%x", h.Sum(nil))
}

//
// Валидации куки с проверкой хеша
//
func (self *CookieReader) parse(value string) string {
    if len(value) == 0 {
        return ""
    }
    segments := strings.Split(value, "+")
    if len(segments) == 1 {
        return ""
    }
    hash := segments[0]
    value = strings.Join(segments[1:], "+")
    if hash == self.hash(value) {
        return value
    }
    return ""
}
