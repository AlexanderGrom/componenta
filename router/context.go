package router

import (
	"net/http"
)

// Контекст запроса
// Содержит расширенные Request и Response
type Ctx struct {
	Req *Request
	Res *Response
}

func NewCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	req := NewRequest(r)
	res := NewResponse(w, r)
	return &Ctx{
		req, res,
	}
}
