package router

import (
	"io"
	"net/http"
	"path"
)

type Multiplexer struct {
	routes Routes
	logger Logger
}

func NewMultiplexer(w io.Writer) *Multiplexer {
	return &Multiplexer{
		routes: NewRoutes(),
		logger: NewLogger(w),
	}
}

func (self *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewCtx(w, r)
	route, err := self.routing(ctx.Req)
	if err != nil {
		urlpath := ctx.Req.URL.Path
		if urlpath[len(urlpath)-1] != '/' {
			ext := path.Ext(urlpath)
			if len(ext) == 0 {
				ctx.Req.URL.Path += "/"
				_, err := self.routing(ctx.Req)
				if err == nil {
					ctx.Res.Status(http.StatusFound)
					ctx.Res.Redirect(ctx.Req.URL.String(), http.StatusMovedPermanently)
					return
				}
			}
		}
		ctx.Res.Status(http.StatusNotFound)
		return
	}

	if err := route.FnChain(ctx); err != nil {
		self.logger.Println(err)
	}
}

func (self *Multiplexer) routing(req *Request) (*Route, error) {
	for _, route := range self.routes[req.Method] {
		if self.checkRoute(req, route) {
			return route, nil
		}
	}
	return nil, ErrRouteNotFound
}

func (self *Multiplexer) checkRoute(req *Request, r *Route) bool {
	path := req.URL.Path
	matches := r.RegExp.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		params := make(URLParams)
		for i, name := range r.RegExp.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		req.Params = params
		return true
	}
	return false
}
