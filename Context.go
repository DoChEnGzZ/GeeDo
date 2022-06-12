package Gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer      http.ResponseWriter
	Req         *http.Request
	Method      string
	Path        string
	StatusCode  int
	Params      map[string]string
	Index       int //middleware index,default -1
	Middlewares []HandlerFunc
	engine      *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
		Index:  -1,
	}
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "1.txt/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTMl(code int, name string,data interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "1.txt/html")
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(http.StatusInternalServerError,err.Error())
	}
}

func (c *Context) Param(key string) interface{} {
	return c.Params[key]
}

func (c *Context) Next() {
	c.Index++
	l := len(c.Middlewares)
	for ; c.Index < l; c.Index++ {
		c.Middlewares[c.Index](c)
	}
}
