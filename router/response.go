package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Обертка над http.ResponseWriter
type Response struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Cookies  *CookieWriter
	Status   int
	data     []byte
	redirect string
}

func NewResponse(w http.ResponseWriter, r *http.Request) *Response {
	c := NewCookieWriter(w)
	return &Response{
		w, r, c, 200, []byte{}, "",
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
	self.data = data
	return nil
}

func (self *Response) Redirect(url string) error {
	self.redirect = url
	return nil
}

func (self *Response) flush() {
	if len(self.redirect) > 0 {
		http.Redirect(self.Writer, self.Request, self.redirect, self.Status)
	} else {
		self.Writer.WriteHeader(self.Status)
		self.Writer.Write(self.data)
	}
}
