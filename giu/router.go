package giu

import "strings"

var (
	httpMethod = []string{"GET", "POST", "PUT", "DELETE", "TRACE", "OPTIONAL"}
)

type HandlersChain []HandlerFunc

type router struct {
	roots    map[string]*node
	handlers map[string]HandlersChain
}

func newRouter() *router {

	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlersChain),
	}
}

func (r *router) addRouter(method string, pattern string, handlerFunc ...HandlerFunc) {

	CheckParam(method, pattern)
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	root, ok := r.roots[method]
	if !ok {
		root = new(node)
		root.fullPath = "/"
		r.roots[method] = root
	}

	root.insertChild(pattern, parts, 0)

	r.handlers[key] = handlerFunc

}

func (r *router) handler(c *Context) {

	path := c.request.URL.Path
	method := c.Method
	n, params := r.getRouter(method, path)
	if n == nil {
		panic("router cannot found")
	}

	c.Params = params

	key := method + "-" + n.fullPath

	handlerFunc, ok := r.handlers[key]
	if !ok {
		panic("get function failed")
	}
	c.handlers = append(c.handlers, handlerFunc...)
	c.Next()

}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {

	parts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	targetNode := root.search(parts, 0)

	if targetNode != nil {
		patternParts := parsePattern(targetNode.fullPath)

		for index, part := range patternParts {

			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts, "/")
			}

		}
		return targetNode, params

	}

	return nil, nil

}
