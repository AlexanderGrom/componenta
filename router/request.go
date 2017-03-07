package router

import (
	"context"
	"net/http"
)

// Обертка над http.Request
type Request struct {
	*http.Request
	Cookies *CookieReader
	Params  URLParams
}

func NewRequest(r *http.Request) *Request {
	c := NewCookieReader(r)
	p := NewURLParams()
	return &Request{
		r, c, p,
	}
}

func (self *Request) WithContext(c context.Context) *Request {
	self.Request = self.Request.WithContext(c)
	return self
}
