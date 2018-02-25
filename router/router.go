package router

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

const (
	GET    = "GET"
	PUT    = "PUT"
	HEAD   = "HEAD"
	POST   = "POST"
	DELETE = "DELETE"
)

var (
	ErrRouteNotFound = errors.New("Route not found")
)

var regexpPlaceholder = regexp.MustCompile(`:([\w]+)`)

type Handler func(*Ctx) error

func (self Handler) apply(ctx *Ctx, fns []appliable, index int) error {
	return self(ctx)
}

func WrapHandler(h http.Handler) Handler {
	return func(ctx *Ctx) error {
		h.ServeHTTP(ctx.Res.Writer, ctx.Req.Request)
		return nil
	}
}

type Route struct {
	*Interceptor
	Method  string
	Pattern string
	Handler Handler
	FnChain func(ctx *Ctx) error
	RegExp  *regexp.Regexp
}

func NewRoute(method string, pattern string, handler Handler) *Route {
	pattern = regexp.QuoteMeta(pattern)
	pattern = regexpPlaceholder.ReplaceAllString(pattern, `(?P<$1>[0-9A-Za-z\-]+)`)
	rexp := regexp.MustCompile(pattern)
	return &Route{
		Interceptor: NewInterceptor(),
		Method:      method,
		Pattern:     pattern,
		Handler:     handler,
		RegExp:      rexp,
	}
}

type Routes map[string][]*Route

func NewRoutes() Routes {
	return Routes{
		GET:    []*Route{},
		PUT:    []*Route{},
		HEAD:   []*Route{},
		POST:   []*Route{},
		DELETE: []*Route{},
	}
}

type Router struct {
	*Interceptor
	*Grouper
	Mux    *Multiplexer
	groups []*Grouper
}

func New(w io.Writer) *Router {
	return &Router{
		Interceptor: NewInterceptor(),
		Grouper:     NewGrouper(""),
		Mux:         NewMultiplexer(w),
	}
}

func (self *Router) Group(prefix string) *Grouper {
	g := NewGrouper(prefix)
	self.groups = append(self.groups, g)
	return g
}

func (self *Router) Handler() http.Handler {
	for method, routes := range self.Routes {
		for _, route := range routes {

			route.FnChain = compose(merge(
				self.middlewares,
				route.middlewares,
				[]appliable{route.Handler},
			))

			self.Mux.routes[method] = append(self.Mux.routes[method], route)
		}
	}

	for _, group := range self.groups {
		for method, routes := range group.Routes {
			for _, route := range routes {

				route.FnChain = compose(merge(
					self.middlewares,
					group.middlewares,
					route.middlewares,
					[]appliable{route.Handler},
				))

				self.Mux.routes[method] = append(self.Mux.routes[method], route)
			}
		}
	}

	return self.Mux
}

type Grouper struct {
	*Interceptor
	Routes Routes
	prefix string
}

func NewGrouper(prefix string) *Grouper {
	return &Grouper{
		Interceptor: NewInterceptor(),
		Routes:      NewRoutes(),
		prefix:      prefix,
	}
}

func (self *Grouper) Get(pattern string, fn Handler) *Route {
	return self.registr(GET, pattern, fn)
}

func (self *Grouper) Post(pattern string, fn Handler) *Route {
	return self.registr(POST, pattern, fn)
}

func (self *Grouper) Head(pattern string, fn Handler) *Route {
	return self.registr(HEAD, pattern, fn)
}

func (self *Grouper) Put(pattern string, fn Handler) *Route {
	return self.registr(PUT, pattern, fn)
}

func (self *Grouper) Delete(pattern string, fn Handler) *Route {
	return self.registr(DELETE, pattern, fn)
}

func (self *Grouper) registr(method string, pattern string, fn Handler) *Route {
	r := NewRoute(method, self.prefix+pattern, fn)
	self.Routes[method] = append(self.Routes[method], r)
	return r
}
