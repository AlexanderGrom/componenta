package router

type Next func()

// Интерфес для Middleware и Handler
// Эти функции будут вызываться используя метод apply
type appliable interface {
	apply(ctx *Ctx, fns []appliable, current int)
}

type Interceptor struct {
	middlewares []appliable
}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

// Добавляет middleware
func (self *Interceptor) Use(fns ...Middleware) *Interceptor {
	a := make([]appliable, len(fns))
	for i, fn := range fns {
		a[i] = fn
	}
	self.middlewares = append(self.middlewares, a...)
	return self
}

func compose(fns []appliable) func(*Ctx) {
	return func(ctx *Ctx) {
		fns[0].apply(ctx, fns, 0)
	}
}

func merge(appliabels ...[]appliable) []appliable {
	all := []appliable{}
	for _, app := range appliabels {
		all = append(all, app...)
	}
	return all
}

type Middleware func(*Ctx, Next)

func (self Middleware) apply(ctx *Ctx, fns []appliable, index int) {
	self(ctx, func() {
		index++
		if len(fns) > index {
			fns[index].apply(ctx, fns, index)
		}
	})
}
