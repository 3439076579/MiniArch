package giu

import (
	"net/http"
	"strings"
	"sync"
)

type HandlerFunc func(c *Context)
type HandlerChain []HandlerFunc

type any = interface{}

type H map[string]interface{}

type Engine struct {
	*RouterGroup

	router *router

	allGroup []*RouterGroup

	pool sync.Pool
}

func New() *Engine {
	e := &Engine{router: newRouter()}
	e.RouterGroup = &RouterGroup{
		engine: e,
	}
	e.allGroup = make([]*RouterGroup, 0)

	e.pool = sync.Pool{
		New: func() interface{} {
			return new(Context)
		},
	}

	return e

}

func (r *Engine) ServeHTTP(w http.ResponseWriter, request *http.Request) {

	var middleware []HandlerFunc

	ctx := r.pool.Get().(*Context)
	ctx.newContext(w, request)

	for _, group := range r.allGroup {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middleware = append(middleware, group.middleware...)
		}
	}

	ctx.handlers = middleware
	r.router.handler(ctx)

	ctx.Clear()

	r.pool.Put(ctx)

}

func (r *Engine) addRouter(method, pattern string, handlerFunc ...HandlerFunc) {
	r.router.addRouter(method, pattern, handlerFunc[0])
}

func (r *Engine) GET(pattern string, handlerFunc ...HandlerFunc) {
	r.addRouter("GET", pattern, handlerFunc...)
}

func (r *Engine) POST(pattern string, handlerFunc ...HandlerFunc) {
	r.addRouter("POST", pattern, handlerFunc...)
}

func (r *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

type RouterGroup struct {
	engine *Engine

	prefix     string
	middleware []HandlerFunc

	parent *RouterGroup
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {

	if !strings.HasPrefix(prefix, "/") {
		panic("prefix must start with slash")
	}

	newGroup := &RouterGroup{
		parent: group,
		prefix: group.parent.prefix + prefix,
		engine: group.engine,
	}

	group.engine.allGroup = append(group.engine.allGroup, newGroup)

	return newGroup

}

func (group *RouterGroup) addRoute(method, pattern string, handlerFunc ...HandlerFunc) {

	pattern = group.prefix + pattern

	group.engine.router.addRouter(method, pattern, handlerFunc...)

}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handlerFunc ...HandlerFunc) {
	group.addRoute("GET", pattern, handlerFunc...)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handlerFunc ...HandlerFunc) {
	group.addRoute("POST", pattern, handlerFunc...)
}

func (group *RouterGroup) Use(handlerFunc ...HandlerFunc) {
	group.middleware = append(group.middleware, handlerFunc...)
}
