package router

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	r := New(nil)
	r.Get("/test/:name", func(ctx *Ctx) error {
		return ctx.Res.Text(ctx.Req.Params.Get("name"))
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
	r := New(nil)
	r.Head("/head", func(ctx *Ctx) error {
		p = true
		return nil
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
	r := New(nil)
	r.Get("/", func(ctx *Ctx) error {
		token := ctx.Req.Context().Value("app.auth.token")
		user := ctx.Req.Context().Value("app.auth.user")
		if token != "123456" || user != "Alexander" {
			return ctx.Res.Status(401)
		}
		return ctx.Res.Text("main")
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
	r := New(nil)

	r.Use(func(ctx *Ctx, next Next) error {
		cxt := context.WithValue(ctx.Req.Context(), "global", true)
		ctx.Req.WithContext(cxt)
		return next()
	})

	g := r.Group("/group")
	g.Use(func(ctx *Ctx, next Next) error {
		cxt := context.WithValue(ctx.Req.Context(), "group1", true)
		ctx.Req.WithContext(cxt)
		return next()
	})
	{
		g.Get("/news", func(ctx *Ctx) error {
			global := ctx.Req.Context().Value("global")
			group1 := ctx.Req.Context().Value("group1")
			route1 := ctx.Req.Context().Value("route1")
			route2 := ctx.Req.Context().Value("route2")
			if global != true || group1 != true || route1 != true || route2 != true {
				return ctx.Res.Status(400)
			}
			return ctx.Res.Text("news")
		}).Use(func(ctx *Ctx, next Next) error {
			cxt := context.WithValue(ctx.Req.Context(), "route1", true)
			ctx.Req.WithContext(cxt)
			return next()
		}).Use(func(ctx *Ctx, next Next) error {
			cxt := context.WithValue(ctx.Req.Context(), "route2", true)
			ctx.Req.WithContext(cxt)
			return next()
		})
	}

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/group/news", nil)
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
	r := New(nil)
	g := r.Group("/group")
	{
		g.Get("/path", func(ctx *Ctx) error {
			return ctx.Res.Text("group/path")
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
	r := New(nil)
	r.Get("/path", func(ctx *Ctx) error {
		ctx.Res.Cookies.Set("userid", "1", 100500)
		ctx.Res.Text("path")
		return nil
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

func TestLogger(t *testing.T) {
	logger := &bytes.Buffer{}
	r := New(logger)
	r.Get("/path", func(ctx *Ctx) error {
		ctx.Res.Text("path")
		return errors.New("error path route")
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/path", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	expected := `path`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}

	log := []byte("error path route")
	if !bytes.Contains(logger.Bytes(), log) {
		t.Errorf("logger containt unexpected text: got %v want %v", logger.String(), string(log))
	}
}

func TestLogger2(t *testing.T) {
	logger := &bytes.Buffer{}
	r := New(logger)
	r.Use(func(ctx *Ctx, next Next) error {
		err := next()
		if err != nil {
			return fmt.Errorf("%s5", err)
		}
		return nil
	})

	g := r.Group("/group")
	g.Use(func(ctx *Ctx, next Next) error {
		err := next()
		if err != nil {
			return fmt.Errorf("%s4", err)
		}
		return nil
	})
	{
		g.Get("/news", func(ctx *Ctx) error {
			return errors.New("1")
		}).Use(func(ctx *Ctx, next Next) error {
			err := next()
			if err != nil {
				return fmt.Errorf("%s3", err)
			}
			return nil
		}).Use(func(ctx *Ctx, next Next) error {
			err := next()
			if err != nil {
				return fmt.Errorf("%s2", err)
			}
			return nil
		})
	}

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/group/news", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	log := []byte("12345")
	if !bytes.Contains(logger.Bytes(), log) {
		t.Errorf("logger containt unexpected text: got %v want %v", logger.String(), string(log))
	}
}

func TestPanicLogger(t *testing.T) {
	logger := &bytes.Buffer{}
	r := New(logger)
	r.Use(func(ctx *Ctx, next Next) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				ctx.Res.Status(http.StatusInternalServerError)
			}
		}()
		return next()
	})
	r.Get("/path", func(ctx *Ctx) error {
		panic("panic path route")
		return nil
	})

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/path", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	log := []byte("panic path route")
	if !bytes.Contains(logger.Bytes(), log) {
		t.Errorf("logger containt unexpected text: got %v want %v", logger.String(), string(log))
	}
}

func TestWrapHandler(t *testing.T) {
	r := New(nil)
	r.Get("/test", WrapHandler(http.NotFoundHandler()))

	mux := r.Handler()

	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `404 page not found`
	if !strings.Contains(res.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}
