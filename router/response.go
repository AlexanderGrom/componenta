package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Обертка над http.ResponseWriter
type Response struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Cookies *CookieWriter
}

func NewResponse(w http.ResponseWriter, r *http.Request) *Response {
	c := NewCookieWriter(w)
	return &Response{
		w, r, c,
	}
}

func (self *Response) Header() http.Header {
	return self.Writer.Header()
}

func (self *Response) Text(text string) error {
	return self.Raw([]byte(text))
}

func (self *Response) Json(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	self.Writer.Header().Set("Content-Type", "application/json")
	return self.Raw(res)
}

func (self *Response) Xml(obj interface{}) error {
	res, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	self.Writer.Header().Set("Content-Type", "application/xml")
	return self.Raw(res)
}

func (self *Response) Raw(data []byte) error {
	_, err := self.Writer.Write(data)
	return err
}

func (self *Response) Redirect(url string, code int) error {
	http.Redirect(self.Writer, self.Request, url, code)
	return nil
}

func (self *Response) Status(code int) error {
	self.Writer.WriteHeader(code)
	return nil
}
