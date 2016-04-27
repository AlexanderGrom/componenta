package router

import (
	"errors"
	"net/http"
	"regexp"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var (
	ErrRouteNotFound = errors.New("Route not found")
)

var regexpPlaceholder = regexp.MustCompile(`:([\w]+)`)

type Handler func(*Ctx) (int, error)

func (self Handler) apply(ctx *Ctx, fns []appliable, index int) {
	status, err := self(ctx)

	if err != nil {
		ctx.Res.Status = http.StatusInternalServerError
	} else {
		ctx.Res.Status = status
	}

	index++
	if len(fns) > index {
		fns[index].apply(ctx, fns, index)
	}
}

type Route struct {
	*interceptor
	Method  string
	Pattern string
	Handler Handler
	FnChain func(ctx *Ctx)
	RegExp  *regexp.Regexp
}

func NewRoute(method string, pattern string, handler Handler) *Route {
	pattern = regexp.QuoteMeta(pattern)
	pattern = regexpPlaceholder.ReplaceAllString(pattern, `(?P<$1>[0-9A-Za-z\-]+)`)
	rexp := regexp.MustCompile(pattern)
	return &Route{
		interceptor: NewInterceptor(),
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
		POST:   []*Route{},
		PUT:    []*Route{},
		DELETE: []*Route{},
	}
}

type Router struct {
	*interceptor
	*group
	mux    *Multiplexer
	groups []*group
}

func New() *Router {
	return &Router{
		interceptor: NewInterceptor(),
		group:       NewGroup(""),
		mux:         NewMultiplexer(),
	}
}

func (self *Router) Group(prefix string) *group {
	g := NewGroup(prefix)
	self.groups = append(self.groups, g)
	return g
}

func (self *Router) Complete() *Multiplexer {
	for method, routes := range self.routes {
		for _, route := range routes {

			route.FnChain = compose(merge(
				self.middlewares,
				route.middlewares,
				[]appliable{route.Handler},
			))

			self.mux.routes[method] = append(self.mux.routes[method], route)
		}
	}

	for _, group := range self.groups {
		for method, routes := range group.routes {
			for _, route := range routes {

				route.FnChain = compose(merge(
					self.middlewares,
					group.middlewares,
					route.middlewares,
					[]appliable{route.Handler},
				))

				self.mux.routes[method] = append(self.mux.routes[method], route)
			}
		}
	}

	return self.mux
}

type group struct {
	*interceptor
	routes Routes
	prefix string
}

func NewGroup(prefix string) *group {
	return &group{
		interceptor: NewInterceptor(),
		routes:      NewRoutes(),
		prefix:      prefix,
	}
}

func (self *group) Get(pattern string, fn Handler) *Route {
	return self.registr(GET, pattern, fn)
}

func (self *group) Post(pattern string, fn Handler) *Route {
	return self.registr(POST, pattern, fn)
}

func (self *group) Put(pattern string, fn Handler) *Route {
	return self.registr(PUT, pattern, fn)
}

func (self *group) Delete(pattern string, fn Handler) *Route {
	return self.registr(DELETE, pattern, fn)
}

func (self *group) registr(method string, pattern string, fn Handler) *Route {
	r := NewRoute(method, self.prefix+pattern, fn)
	self.routes[method] = append(self.routes[method], r)
	return r
}
