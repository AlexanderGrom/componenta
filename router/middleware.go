package router

type Next func()

// Интерфес для Middleware и Handler
// Эти функции будут вызываться используя метож apply
type appliable interface {
	apply(ctx *Ctx, fns []appliable, current int)
}

type interceptor struct {
	middlewares []appliable
}

func NewInterceptor() *interceptor {
	return &interceptor{}
}

// Добавляет middleware
func (self *interceptor) Use(fns ...Middleware) {
	a := make([]appliable, len(fns))
	for i, fn := range fns {
		a[i] = fn
	}
	self.middlewares = append(self.middlewares, a...)
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
