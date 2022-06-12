package Gee

import "strings"

type Node struct {
	pattern string //当前结点路由，/p/:lang
	part string //该结点路由尾
	children []*Node //子结点
	isWild bool //是否精准匹配,当为* :时为非精准匹配，iswild为ture
}

func (n *Node) matchChild(part string) *Node{
	for _,child:=range n.children{
		if part==child.part||child.isWild{
			return child
		}
	}
	return nil
}

func (n *Node) matchChildren(part string) []*Node{
	children:=make([]*Node,0)
	for _,child:=range n.children{
		if part==child.part||child.isWild{
			children=append(children,child)
		}
	}
	return children
}

func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
