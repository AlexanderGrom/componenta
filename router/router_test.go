package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheck(t *testing.T) {
	r := New()
	r.Get("/test/:name", func(ctx *Ctx) (int, error) {
		ctx.Res.Text(ctx.Req.Params.Get("name"))
		return 200, nil
	})

	mux := r.Complete()

	req, err := http.NewRequest("GET", "/test/check", nil)
	if err != nil {
		t.Fatal(err)
	}

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

func TestContext(t *testing.T) {
	r := New()
	r.Get("/", func(ctx *Ctx) (int, error) {
		token := ctx.Req.Context().Value("app.auth.token")
		user := ctx.Req.Context().Value("app.auth.user")
		if token != "123456" || user != "Alexander" {
			return 401, nil
		}
		ctx.Res.Text("main")
		return 200, nil
	})

	mux := r.Complete()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	r.Get("/news", func(ctx *Ctx) (int, error) {
		global := ctx.Req.Context().Value("global")
		route := ctx.Req.Context().Value("route")
		if global != true || route != true {
			return 400, nil
		}
		ctx.Res.Text("news")
		return 200, nil
	}).Use(func(ctx *Ctx, next Next) {
		cxt := context.WithValue(ctx.Req.Context(), "route", true)
		ctx.Req.WithContext(cxt)
		next()
	})

	mux := r.Complete()

	req, err := http.NewRequest("GET", "/news", nil)
	if err != nil {
		t.Fatal(err)
	}

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
		g.Get("/path", func(ctx *Ctx) (int, error) {
			ctx.Res.Text("group/path")
			return 200, nil
		})
	}

	mux := r.Complete()

	req, err := http.NewRequest("GET", "/group/path", nil)
	if err != nil {
		t.Fatal(err)
	}

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
	r.Get("/path", func(ctx *Ctx) (int, error) {
		ctx.Res.Cookies.Set("userid", "1", 100500)
		ctx.Res.Text("path")
		return 200, nil
	})

	mux := r.Complete()

	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}

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
