
## Componenta / Router

Простой роутер...

```go
package main

import (
    "github.com/AlexanderGrom/componenta/router"
    "log"
    "net/http"
)

func main() {  
    r := router.New(nil)

    r.Get("/", func(ctx *router.Ctx) error {
        ctx.Res.Cookies.Set("test", "Home Page", 100500)
        ctx.Res.Text("Hello World")
        return nil
    })

    r.Get("/test", func(ctx *router.Ctx) error {
        test := ctx.Req.Cookies.Get("test")
        return ctx.Res.Text("Cookie Value: " + test)
    })

    r.Get("/test/:name", func(ctx *router.Ctx) error {
        return ctx.Res.Text(ctx.Req.Params.Get("name"))
    })

    r.Get("/name", func(ctx *router.Ctx) error {
        return ctx.Res.Redirect("/test/name", 301)
    })

    r.Use(func(ctx *router.Ctx, next router.Next) error {
        log.Println("Global Middleware")
        return next()
    })

    g := r.Group("/group")
    g.Use(func(ctx *router.Ctx, next router.Next) error {
        log.Println("Group Middleware")
        return next()
    })
    {
        g.Get("/path", func(ctx *router.Ctx) error {
            return ctx.Res.Text("Address: /group/path")
        }).Use(func(ctx *router.Ctx, next router.Next) error {
            log.Println("Route Middleware")
            return next()
        })
    }

    if err := http.ListenAndServe(":8080", r.Handler()); err != nil {
        log.Fatalln("ListenAndServe:", err)
    }
}
```
