package router

import (
    "errors"
    "log"
    "net/http"
    "path"
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

type Route struct {
    Method  string
    Pattern string
    Handler Handler
    RegExp  *regexp.Regexp
}

func NewRoute(method string, pattern string, handler Handler) *Route {
    pattern = regexp.QuoteMeta(pattern)
    pattern = regexpPlaceholder.ReplaceAllString(pattern, `(?P<$1>[0-9A-Za-z\-]+)`)
    rexp := regexp.MustCompile(pattern)
    return &Route{
        method, pattern, handler, rexp,
    }
}

type Router struct {
    routes map[string][]*Route
}

func New() *Router {
    return &Router{
        routes: map[string][]*Route{
            GET:    []*Route{},
            POST:   []*Route{},
            PUT:    []*Route{},
            DELETE: []*Route{},
        },
    }
}

func (self *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := NewCtx(w, r)

    defer ctx.Res.flush()

    route, err := self.routing(ctx.Req)

    if err != nil {
        urlpath := ctx.Req.URL.Path
        if urlpath[len(urlpath)-1] != '/' {
            ext := path.Ext(urlpath)
            if len(ext) == 0 {
                ctx.Req.URL.Path += "/"
                _, err := self.routing(ctx.Req)
                if err == nil {
                    ctx.Res.Status = http.StatusFound
                    ctx.Res.Redirect(ctx.Req.URL.String())
                    return
                }
            }
        }
        ctx.Res.Status = http.StatusNotFound
        return
    }

    status, err := route.Handler(ctx)
    if err != nil {
        ctx.Res.Status = http.StatusInternalServerError
    } else {
        ctx.Res.Status = status
    }
}

func (self *Router) Get(pattern string, fn Handler) {
    self.registr(NewRoute(GET, pattern, fn))
}

func (self *Router) Post(pattern string, fn Handler) {
    self.registr(NewRoute(POST, pattern, fn))
}

func (self *Router) Put(pattern string, fn Handler) {
    self.registr(NewRoute(PUT, pattern, fn))
}

func (self *Router) Delete(pattern string, fn Handler) {
    self.registr(NewRoute(DELETE, pattern, fn))
}

func (self *Router) registr(r *Route) {
    self.routes[r.Method] = append(self.routes[r.Method], r)
}

func (self *Router) routing(req *Request) (*Route, error) {
    for _, route := range self.routes[req.Method] {
        if self.checkRoute(req, route) {
            return route, nil
        }
    }
    return nil, ErrRouteNotFound
}

func (self *Router) checkRoute(req *Request, r *Route) bool {
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

