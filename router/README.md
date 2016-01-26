
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

    r.Get("/", func(ctx *router.Ctx) (int, error) {
        ctx.Res.Cookies.Set("test", "Home Page", 100500)
        ctx.Res.Text("Hello World")
        return 200, nil
    })

    r.Get("/test", func(ctx *router.Ctx) (int, error) {
        test := ctx.Req.Cookies.Get("test")
        ctx.Res.Text("Cookie Value: " + test)
        return 200, nil
    })

    r.Get("/test/:name", func(ctx *router.Ctx) (int, error) {
        ctx.Res.Text(ctx.Req.Params.Get("name"))
        return 200, nil
    })

    r.Get("/name", func(ctx *router.Ctx) (int, error) {
        ctx.Res.Redirect("/test/name")
        return 301, nil
    })

    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalln("ListenAndServe:", err)
    }
}
```