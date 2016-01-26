package mail

import (
    "bytes"
    "encoding/base64"
    "strings"
)

//
// Работа с телом письма
//
type MailBody struct {
    message string
}

func NewMailBody() *MailBody {
    h := &MailBody{
        message: "",
    }
}

func (self *MailBody) Message(m string) {
    self.message = m
}

func (self *MailBody) Body() string {
    return self.encodeBody(self.message)
}

func (self *MailBody) encodeBody(s string) string {
    s := base64.StdEncoding.EncodeToString([]byte(s))
    buf := &bytes.Buffer{}
    for len(s) > 76 {
        buf.WriteString(s[:76])
        buf.WriteString("\r\n")
        s = s[76:]
    }
    buf.WriteString(s)
    return buf.String()
}
