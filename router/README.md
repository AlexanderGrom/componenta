
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
    r := router.New()

    r.Get("/", func(ctx *router.Ctx) int {
        ctx.Res.Cookies.Set("test", "Home Page", 100500)
        ctx.Res.Text("Hello World")
        return 200
    })

    r.Get("/test", func(ctx *router.Ctx) int {
        test := ctx.Req.Cookies.Get("test")
        ctx.Res.Text("Cookie Value: " + test)
        return 200
    })

    r.Get("/test/:name", func(ctx *router.Ctx) int {
        ctx.Res.Text(ctx.Req.Params.Get("name"))
        return 200
    })

    r.Get("/name", func(ctx *router.Ctx) int {
        ctx.Res.Redirect("/test/name")
        return 301
    })

    r.Use(func(ctx *router.Ctx, next router.Next) {
		log.Println("Global Middleware")
		next()
	})

	g := r.Group("/group")
	g.Use(func(ctx *router.Ctx, next router.Next) {
		log.Println("Group Middleware")
		next()
	})
	{
		g.Get("/path", func(ctx *router.Ctx) int {
			ctx.Res.Text("Address: /group/path")
			return 200
		}).Use(func(ctx *router.Ctx, next router.Next) {
			log.Println("Route Middleware")
			next()
		})
	}

    mux := r.Complete()

    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalln("ListenAndServe:", err)
    }
}
```
