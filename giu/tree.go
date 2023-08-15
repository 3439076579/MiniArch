package giu

import (
	"unsafe"
)

func NewMethodTree() trees {
	return make([]*methodTree, 0)
}

type methodTree struct {
	method string
	root   *node
}
type trees []*methodTree

func (t trees) GetRoot(method string) *node {
	for _, tree := range t {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

func longestPrefix(p, s string) int {

	min := func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	var i int
	for i = 0; i < min(len(p), len(s)); i++ {
		if p[i] != s[i] {
			break
		}
	}
	return i
}

type nodeType uint8

const (
	root nodeType = iota
	static
	param
	catchAll
)

type node struct {
	path       string
	parent     *node
	paramChild *node
	anyChild   *node
	indices    string
	children   []*node
	nType      nodeType
	handlers   HandlerChain
}

func newNode() *node {
	return &node{
		children: make([]*node, 0),
	}
}

func (n *node) addRoute(path string, handlers HandlerChain) {
	// Empty tree
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, handlers)
		n.nType = root
		return
	}
walk:
	i := longestPrefix(path, n.path)

	if i < len(n.path) {
		child := &node{
			path:       n.path[i:],
			parent:     n,
			children:   n.children,
			indices:    n.indices,
			anyChild:   n.anyChild,
			paramChild: n.paramChild,
			handlers:   n.handlers,
			nType:      static,
		}

		for _, node := range n.children {
			node.parent = child
		}
		if n.paramChild != nil {
			n.paramChild.parent = child
		}
		if n.anyChild != nil {
			n.anyChild.parent = child
		}

		n.paramChild = nil
		n.anyChild = nil
		n.children = []*node{child}
		n.path = path[:i]
		n.indices = b2s([]byte{child.path[0]})

		if i < len(path) {

			path = path[i:]

			if path[0] != ':' && path[0] != '*' {
				child := &node{
					nType:  static,
					parent: n,
				}
				n.children = append(n.children, child)
				n.indices += b2s([]byte{path[0]})
				n = child
			}
			n.insertChild(path, handlers)

		}

		n.handlers = handlers

		return
	}

	if i < len(path) {
		path = path[i:]

		for i := 0; i < len(n.indices); i++ {
			if path[0] == n.indices[i] {
				n = n.children[i]
				goto walk
			}
		}

		if path[0] != ':' && path[0] != '*' && n.nType != catchAll {
			child := &node{
				parent: n,
				nType:  static,
			}
			n.children = append(n.children, child)
			n.indices += b2s([]byte{path[0]})
			child.insertChild(path, handlers)
		} else {
			if n.nType == catchAll {
				panic("wildcard * must in the end of path")
			}
			switch path[0] {
			case ':':
				if n.paramChild != nil &&
					len(path) >= len(n.paramChild.path) &&
					path[:len(n.paramChild.path)] == n.paramChild.path &&
					(len(path) <= len(n.paramChild.path) ||
						path[len(n.paramChild.path)] == '/') {
					n = n.paramChild
					goto walk
				} else if n.paramChild == nil {
					seg, _, _ := findWildSeg(path)
					if len(path) == len(seg) {
						child := &node{
							path:   seg,
							parent: n,
							nType:  param,
						}
						n.paramChild = child
						n = n.paramChild
						n.handlers = handlers
					} else {
						child := &node{
							path:   seg[:len(seg)-1],
							parent: n,
							nType:  param,
						}
						path = path[len(seg)-1:]
						n.paramChild = child
						c := &node{
							nType:  static,
							parent: child,
						}
						child.children = append(child.children, c)
						n = c
						n.insertChild(path, handlers)
					}
				} else {
					panic("wildcard conflict")
				}
			case '*':
				if n.anyChild != nil {
					panic("")
				}
				child := &node{
					nType:  catchAll,
					parent: n,
				}
				n.anyChild = child
				n = child
				n.insertChild(path, handlers)
			}
		}
	}
}

// InsertChild node must be a empty node
func (n *node) insertChild(path string, handlers HandlerChain) {
	for {
		seg, i, valid := findWildSeg(path)
		if i < 0 {
			goto noWildCard
		}
		if !valid {
			panic("too many wildcard exist in the same path segment," +
				"only one wildcard is allowed in per path segment")
		}
		if len(seg) < 2 {
			panic("wildcard must be named")
		}
		if seg[0] == ':' {
			if i > 0 {
				n.path = path[:i]
			}
			child := &node{
				path:     seg,
				nType:    param,
				children: make([]*node, 0),
				parent:   n,
			}

			n.paramChild = child
			path = path[i+len(seg):]
			n = n.paramChild

			if len(path) != 0 {
				newNode := &node{
					parent: n,
				}
				switch path[0] {
				case ':':
					n.paramChild = newNode
					newNode.nType = param
				case '*':
					n.anyChild = newNode
					newNode.nType = catchAll
				default:
					n.children = append(n.children, newNode)
					newNode.nType = static
					n.indices += b2s([]byte{path[0]})
				}
				n = newNode
				continue
			}
			n.handlers = handlers
			return
		}

		if i+len(seg) != len(path) {
			panic("wildcard * must be the end of path")
		}
		// * wildcard
		child := &node{
			path:   seg,
			parent: n,
			nType:  catchAll,
		}
		n.path = path[:i]
		n.anyChild = child
		n = child
		n.handlers = handlers
		return
	}
noWildCard:
	n.path = path
	n.handlers = handlers

}

func findWildSeg(path string) (string, int, bool) {
	var valid = true
	for start := 0; start < len(path); start++ {
		if path[start] != ':' && path[start] != '*' {
			continue
		}
		if start > 0 && path[start-1] != '/' {
			panic("wildcard must have / before")
		}

		for end := start + 1; end < len(path); end++ {
			switch path[end] {
			case '/':
				return path[start:end], start, valid
			case '*', ':':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, true
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type Param struct {
	Key   string
	Value string
}

type nodeValue struct {
	handlers HandlerChain
	params   []Param
	//searchList []*node
	fullPath string
}

type skippedNode struct {
	validNode  *node
	path       string
	paramCount int64
	fullPath   string
}

func (n *node) Search(path string) (value nodeValue) {
	valid := checkSearchPath(path)
	if !valid {
		panic("invalid search path")
	}
	var skippedNodeList []*skippedNode
	var globalParamCount int64
	var fullPath string
walk:
	for {
		prefix := n.path
		if len(prefix) < len(path) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]
				fullPath += prefix

				c := path[0]
				for i := 0; i < len(n.indices); i++ {
					if c == n.indices[i] {
						if n.paramChild != nil || n.anyChild != nil {
							child := &skippedNode{
								path: prefix + path,
								validNode: &node{
									path:       n.path,
									parent:     n.parent,
									paramChild: n.paramChild,
									anyChild:   n.anyChild,
									children:   n.children,
									nType:      n.nType,
									handlers:   n.handlers,
								},
								paramCount: globalParamCount,
								fullPath:   fullPath,
							}
							skippedNodeList = append(skippedNodeList, child)
						}
						n = n.children[i]
						continue walk
					}
				}

				if n.paramChild == nil && n.anyChild == nil {
					if len(skippedNodeList) > 0 {
						validNode := skippedNodeList[len(skippedNodeList)-1]
						skippedNodeList = skippedNodeList[:len(skippedNodeList)-1]
						value.params = value.params[:validNode.paramCount]
						path = validNode.path
						n = validNode.validNode
						fullPath = validNode.fullPath[:len(validNode.fullPath)-len(n.path)]
						continue walk
					}
				}

				// 走到这里说明有通配符，且无static结点，走通配符
				if n.paramChild != nil {
					n = n.paramChild

					end := 0
					for ; end < len(path) && path[end] != '/'; end++ {
					}

					if value.params == nil {
						value.params = make([]Param, 0)
					}
					paramValue := path[:end]
					fullPath += n.path
					paramKey := n.path[1:]
					value.params = append(value.params, Param{Key: paramKey,
						Value: paramValue})
					globalParamCount++

					// 代表还需要继续走下去
					if end < len(path) {
						path = path[end:]
						if len(n.children) > 0 {
							for i := 0; i < len(n.indices); i++ {
								if n.indices[i] == '/' {
									n = n.children[i]
									continue walk
								}
							}
						}
						if path == "/" && n.handlers != nil {
							value.fullPath = fullPath
							value.handlers = n.handlers
							return
						}
						return
					}
					//走到这里 说明end==len(path)，说明该结尾了
					if n.handlers != nil {
						value.handlers = n.handlers
						value.fullPath = fullPath
						return
					}
					for i := 0; i < len(n.indices); i++ {
						if n.children[i].path == "/" {
							value.handlers = n.children[i].handlers
							value.fullPath = fullPath + "/"
							return
						}
					}

					return
				}

				if n.anyChild != nil {
					n = n.anyChild
					if value.params == nil {
						value.params = make([]Param, 0)
					}
					fullPath += n.path
					paramKey := n.path[1:]
					value.params = append(value.params, Param{
						Key:   paramKey,
						Value: path,
					})
					globalParamCount++

					value.handlers = n.handlers
					value.fullPath = fullPath
					return
				}
			}
		}

		if prefix == path {
			fullPath += prefix
			if n.handlers != nil {
				value.handlers = n.handlers
				value.fullPath = fullPath
				return
			}
			if n.anyChild != nil {
				value.handlers = n.anyChild.handlers
				value.fullPath = fullPath + n.anyChild.path
				return
			}
			if path[len(path)-1] != '/' {
				for i := 0; i < len(n.indices); i++ {
					if n.children[i].path == "/" &&
						n.children[i].handlers != nil {
						value.handlers = n.children[i].handlers
						value.fullPath = fullPath + "/"
						return
					}
					if n.children[i].path == "/" &&
						n.children[i].anyChild != nil &&
						n.children[i].anyChild.handlers != nil {
						value.handlers = n.children[i].anyChild.handlers
						value.fullPath = fullPath + "/" + n.children[i].anyChild.path
						return
					}
				}
			}
			if path == "/" {
				if n.parent.handlers != nil {
					value.handlers = n.parent.handlers
					value.fullPath = fullPath[:len(fullPath)-1]
					return
				}
			}
			if len(skippedNodeList) > 0 {
				validNode := skippedNodeList[len(skippedNodeList)-1]
				skippedNodeList = skippedNodeList[:len(skippedNodeList)-1]
				value.params = value.params[:validNode.paramCount]
				path = validNode.path
				n = validNode.validNode
				fullPath = validNode.fullPath[:len(validNode.fullPath)-len(n.path)]
				continue walk
			}
			return
		}

		if len(path)+1 == len(prefix) &&
			path[len(path)-1] == '/' &&
			path[:len(path)-1] == prefix &&
			n.handlers != nil {
			value.handlers = n.handlers
			value.fullPath += prefix
			return
		}

		if path == "/" && n.parent.handlers != nil {
			value.handlers = n.parent.handlers
			value.fullPath = fullPath
			return
		}
		if len(skippedNodeList) > 0 {
			validNode := skippedNodeList[len(skippedNodeList)-1]
			skippedNodeList = skippedNodeList[:len(skippedNodeList)-1]
			value.params = value.params[:validNode.paramCount]
			path = validNode.path
			n = validNode.validNode
			fullPath = validNode.fullPath[:len(validNode.fullPath)-len(n.path)]
			continue walk
		}
	}
}

func checkSearchPath(path string) (valid bool) {
	valid = true
	if path[0] != '/' {
		valid = false
		return
	}
	var index int
	for index < len(path) {
		switch path[index] {
		case '/':
			if (index < len(path)-1) && path[index+1] == '/' {
				valid = false
				return
			}
			index++
		case ':', '*':
			valid = false
			return
		default:
			var end int
			for ; index+end < len(path) && path[index+end] != '/'; end++ {
			}
			index += end
		}
	}
	return
}
