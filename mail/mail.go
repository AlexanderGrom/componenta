package mail

import (
    "bytes"
    "errors"
    "net/smtp"
)

type Mail struct {
    from    string
    to      string
    headers *MailHeaders
    body    *MailBody
}

func New() *Mail {
    return &Mail{
        headers: NewMailHeaders(),
        body:    NewMailBody(),
    }
}

//
// Адрес и имя получателя
//
func (self *Mail) To(address string, name string) {
    self.headers.To(address, name)
    self.to = address
}

//
// Адрес и имя отправителя
//
func (self *Mail) From(address string) {
    self.headers.From(address)
    self.from = address
}

//
// Заголовок письма
//
func (self *Mail) Subject(subject string) {
    self.headers.Subject(subject)
}

//
// Текст письма
//
func (self *Mail) Text(subject string) {
    self.body.Message(subject)
}

//
// Отправка
//
func (self *Mail) Send() error {
    headers := self.headers.Headers()

    for _, v := range []string{"From", "To", "Subject"} {
        if _, ok := headers[v]; !ok {
            return errors.New("Mail: not \"" + v + "\" header")
        }
    }

    buf := &bytes.Buffer{}
    for k, v := range headers {
        buf.WriteString(k)
        buf.WriteString(": ")
        buf.WriteString(v)
        buf.WriteString("\r\n")
    }
    buf.WriteString("\r\n")
    buf.WriteString(self.body.Body())

    msg := buf.Bytes()

    err = smtp.SendMail("localhost:25", nil, self.from, []string{self.to}, msg)
    
    if err != nil {
        return err
    }
}
