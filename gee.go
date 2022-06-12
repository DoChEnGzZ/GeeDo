package Gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	//groups []*RouterGroup
	htmlTemplate *template.Template
	funcMap template.FuncMap
}

type RouterGroup struct{
	r *Router
	prefix string
	middlewares []HandlerFunc
	groups []*RouterGroup
	parent *RouterGroup
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c:=NewContext(writer,request)
	c.engine=e
	//c.Middlewares=append(c.Middlewares, e.middlewares...)
	middlewares :=make([]HandlerFunc,0)
	//加入其父路由组的中间件
	middlewares =append(middlewares, e.middlewares...)
	for _,group:=range e.groups{
		if strings.HasPrefix(request.URL.Path,group.prefix){
			//加入对应group的中间件
			middlewares=append(middlewares, group.middlewares...)
		}
	}
	c.Middlewares= middlewares
	e.r.handle(c)
}

func New() *Engine {
	engine:=&Engine{
	}
	engine.RouterGroup=&RouterGroup{
		r:NewRouter(),
	}
	engine.groups=[]*RouterGroup{}
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap){
	e.funcMap=funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func Default() *Engine {
	engine:=New()
	engine.Use(Logger(),Recovery())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	newGroup:=&RouterGroup{
		r:group.r,
		prefix: group.prefix+prefix,
		parent: group,
	}
	group.groups=append(group.groups, newGroup)
	return newGroup
}

//method get pattern path handler func
func (group *RouterGroup) addRoute(method string,cmp string,handler HandlerFunc)  {
	pattern:= group.prefix+cmp
	log.Printf("Route %4s - %s", method, pattern)
	group.r.addRoute(method,pattern,handler)
}

func (group *RouterGroup) Get(pattern string,handler HandlerFunc)  {
	group.addRoute("GET",pattern,handler)
}

func (group *RouterGroup) Post(pattern string,handler HandlerFunc)  {
	group.addRoute("POST",pattern,handler)
}

func (group *RouterGroup) Use(Middlewares...HandlerFunc)  {
	group.middlewares=append(group.middlewares, Middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string,system http.FileSystem)HandlerFunc{
	absolutePath:=path.Join(group.prefix, relativePath)
	fs:=http.StripPrefix(absolutePath,http.FileServer(system))
	return func(c *Context) {
		file:=c.Param("filepath")
		if _,err:=system.Open(file.(string));err!=nil{
			c.String(http.StatusNotFound,err.Error())
			return
		}
		fs.ServeHTTP(c.Writer,c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string,root string)  {
	handler:=group.createStaticHandler(relativePath,http.Dir(root))
	urlPattern:=path.Join(relativePath,"/*filepath")
	group.Get(urlPattern,handler)
}

func (e *Engine) Run(addr string) (err error){
	return http.ListenAndServe(addr,e)
}
