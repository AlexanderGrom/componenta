package mail

import (
    "mime"
    "regexp"
    "strings"
)

//
// Работа с заголовками письма
//
type MailHeaders struct {
    headers map[string]string
}

func NewMailHeaders() *MailHeaders {
    h := &MailHeaders{
        headers: make(map[string]string),
    }
    h.headers["MIME-Version"] = "1.0"
    h.headers["X-Mailer"] = "Go Mailer"
    h.headers["Content-Type"] = "text/plain; charset=utf-8"
    h.headers["Content-Transfer-Encoding"] = "base64"
}

//
// Добавление заголовка получателя
//
func (self *MailHeaders) To(address string) {
    if !self.isEmail(address) {
        return
    }
    self.headers["To"] = self.secureHeader(address)
}

//
// Добавление заголовка отправителя
//
func (self *MailHeaders) From(address string, name string) {
    if !self.isEmail(address) {
        return
    }
    self.headers["From"] = self.encodeHeader(name) + " <" + self.secureHeader(address) + ">"
    self.headers["Reply-To"] = self.headers["From"]
    self.headers["Return-Path"] = self.headers["From"]
}

//
// Добавление заголовка темы письма
//
func (self *MailHeaders) Subject(subject string) {
    self.headers["Subject"] = self.encodeHeader(subject)
}

//
// Отдача имеющихся заголовков
//
func (self *MailHeaders) Headers() map[string][]string {
    return self.headers
}

//
// Очистка
//
func (self *MailHeaders) secureHeader(s string) string {
    s = strings.Replace(s, "\r", "", -1)
    s = strings.Replace(s, "\n", "", -1)
    return strings.TrimSpace(s)
}

//
// Кодирует строки по стандарту RFC 2045
//
func (self *MailHeaders) encodeHeader(s string) string {
    return mime.BEncoding.Encode("UTF-8", self.secureHeader(s))
}

//
// Проверка на валидный e-mail адрес
//
func (self *MailHeaders) isEmail(s string) bool {
    if ok, _ := MatchString(`^[\d\p{L}\.\_\+\-]+@[\d\p{L}\-]+(\.[\d\p{L}\-]+)+$`, s); !ok {
        return false
    }
    return true
}
