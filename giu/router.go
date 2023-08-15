package giu

import (
	"regexp"
)

var (
	httpMethod    = []string{"GET", "POST", "PUT", "DELETE", "TRACE", "OPTIONAL"}
	httpMethodReg = `[A-Z]{3,8}`
)

type router struct {
	methodTree trees
}

func newRouter() *router {
	return &router{
		methodTree: NewMethodTree(),
	}
}

func (r *router) addRouter(method string, pattern string, handlerFunc ...HandlerFunc) {
	if pattern[0] != '/' {
		panic("path must start with '/'")
	}
	if method == "" {
		panic("method name is not allowed to be empty")
	}
	if len(handlerFunc) == 0 {
		panic("null handlers")
	}
	root := r.methodTree.GetRoot(method)
	if root == nil {
		root = newNode()
		r.methodTree = append(r.methodTree, &methodTree{
			method: method,
			root:   root,
		})
	}
	root.addRoute(pattern, handlerFunc)

}

func (r *router) handler(c *Context) {

	path := c.request.URL.Path
	method := c.Method

	matched, _ := regexp.MatchString(httpMethodReg, c.Method)
	if !matched {
		panic("http method" + c.Method + "is not valid")
	}

	for _, tree := range r.methodTree {
		if tree.method != method {
			continue
		}
		root := tree.root
		value := root.Search(path)
		if value.params != nil {
			c.Params = value.params
		}
		if value.handlers != nil {
			c.handlers = append(c.handlers, value.handlers...)
			c.FullPath = value.fullPath
			c.handlers[0](c)
			c.Next()
		}
	}

}
