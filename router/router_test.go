package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheck(t *testing.T) {
	r := New()
	r.Get("/test/:name", func(ctx *Ctx) int {
		ctx.Res.Text(ctx.Req.Params.Get("name"))
		return 200
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/test/check", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `check`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestHead(t *testing.T) {
	p := false
	r := New()
	r.Head("/head", func(ctx *Ctx) int {
		p = true
		return 200
	})

	mux := r.Handler()

	req := httptest.NewRequest("HEAD", "/head", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !p {
		t.Errorf("handler don't work")
	}
}

func TestContext(t *testing.T) {
	r := New()
	r.Get("/", func(ctx *Ctx) int {
		token := ctx.Req.Context().Value("app.auth.token")
		user := ctx.Req.Context().Value("app.auth.user")
		if token != "123456" || user != "Alexander" {
			return 401
		}
		ctx.Res.Text("main")
		return 200
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "app.auth.token", "123456")
	ctx = context.WithValue(ctx, "app.auth.user", "Alexander")

	mux.ServeHTTP(res, req.WithContext(ctx))

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `main`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestMiddleware(t *testing.T) {
	r := New()

	r.Use(func(ctx *Ctx, next Next) {
		cxt := context.WithValue(ctx.Req.Context(), "global", true)
		ctx.Req.WithContext(cxt)
		next()
	})

	r.Get("/news", func(ctx *Ctx) int {
		global := ctx.Req.Context().Value("global")
		route1 := ctx.Req.Context().Value("route1")
		route2 := ctx.Req.Context().Value("route2")
		if global != true || route1 != true || route2 != true {
			return 400
		}
		ctx.Res.Text("news")
		return 200
	}).Use(func(ctx *Ctx, next Next) {
		cxt := context.WithValue(ctx.Req.Context(), "route1", true)
		ctx.Req.WithContext(cxt)
		next()
	}).Use(func(ctx *Ctx, next Next) {
		cxt := context.WithValue(ctx.Req.Context(), "route2", true)
		ctx.Req.WithContext(cxt)
		next()
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/news", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `news`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestGroup(t *testing.T) {
	r := New()
	g := r.Group("/group")
	{
		g.Get("/path", func(ctx *Ctx) int {
			ctx.Res.Text("group/path")
			return 200
		})
	}

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/group/path", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `group/path`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestCookies(t *testing.T) {
	r := New()
	r.Get("/path", func(ctx *Ctx) int {
		ctx.Res.Cookies.Set("userid", "1", 100500)
		ctx.Res.Text("path")
		return 200
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/path", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if res.HeaderMap.Get("Set-Cookie")[:7] != "userid=" {
		t.Errorf("cookie not set")
	}

	expected := `path`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}
