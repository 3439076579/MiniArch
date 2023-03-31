package giu

import "fmt"

type node struct {
	part     string
	fullPath string
	children []*node
	isWild   bool
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.fullPath, n.part, n.isWild)
}

func (n *node) matchChild(part string) *node {

	for i := 0; i < len(n.children); i++ {
		if n.children[i].part == part || n.children[i].isWild == true {
			return n.children[i]
		}
	}
	return nil

}

func (n *node) insertChild(pattern string, parts []string, height int) {

	if len(parts) == height {
		if n.fullPath == "" {
			n.fullPath = pattern
		} else {
			panic("router collision")
		}
		return
	}

	part := parts[height]

	// 从子节点集合中找出路由块匹配的，如果找不到返回nil
	child := n.matchChild(part)
	if child == nil {
		newNode := new(node)
		newNode.part = part
		newNode.isWild = part[0] == '*' || part[0] == ':'
		n.children = append(n.children, newNode)
		child = newNode
	}

	child.insertChild(pattern, parts, height+1)

}

func (n *node) searchChildren(part string) []*node {

	children := make([]*node, 0)

	for _, v := range n.children {
		if v.part == part || v.isWild == true {
			children = append(children, v)
		}
	}

	return children

}

func (n *node) search(parts []string, height int) *node {

	if len(parts) == height || parts[height][0] == '*' {
		if n.fullPath == "" {
			return nil
		}
		return n
	}

	// 去寻找多个匹配的Children
	children := n.searchChildren(parts[height])
	for _, child := range children {

		result := child.search(parts, height+1)
		if result != nil {
			return result
		}

	}

	return nil

}
