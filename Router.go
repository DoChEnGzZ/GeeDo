package Gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots map[string]*Node
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots: make(map[string]*Node),
		handlers: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//func (r *Router) addRouter(method string,pattern string,handler HandlerFunc)  {
//	log.Printf("Route %4s - %s", method, pattern)
//	key:=method+"-"+pattern
//	r.handlers[key]=handler
//}


func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &Node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*Node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//执行设置路由时的handler函数
func (r *Router) handle(c *Context){
	n,params:=r.getRoute(c.Method,c.Path)
	if n!=nil {
		c.Params=params
		key:=c.Method+"-"+n.pattern
		c.Middlewares=append(c.Middlewares, r.handlers[key])
		//r.handlers[key](c)
	}else {
		c.Middlewares=append(c.Middlewares, func(c *Context) {
			c.String(http.StatusNotFound,"404 NOT FOUND")
		})
	}
	c.Next()
}