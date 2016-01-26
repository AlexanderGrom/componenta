package router

import (
    "net/http"
)

//
// Обертка над http.Request
//
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
