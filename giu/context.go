package giu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Context struct {
	writer   http.ResponseWriter
	request  *http.Request
	Path     string
	Method   string
	FullPath string

	mu sync.RWMutex

	Keys       map[string]interface{}
	StatusCode int64
	Params     []Param

	handlers []HandlerFunc
	index    int
}

func (c *Context) newContext(w http.ResponseWriter, r *http.Request) {

	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	if c.handlers == nil {
		c.handlers = make([]HandlerFunc, 0)
	}

	c.Method = r.Method
	c.request = r
	c.Path = r.URL.Path
	c.writer = w

}

func (c *Context) Set(key string, value interface{}) {

	c.mu.Lock()
	c.Keys[key] = value
	c.mu.Unlock()

}

func (c *Context) Get(key string) interface{} {

	c.mu.RLock()
	value, ok := c.Keys[key]

	if !ok {
		panic("cannot get value from Keys")
	}

	c.mu.RUnlock()
	return value

}

func (c *Context) Next() {

	c.index++

	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}

}

func (c *Context) PostForm(key string) string {
	return c.request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

func (c *Context) Status(code int64) {
	c.StatusCode = code
	c.writer.WriteHeader(int(code))
}

func (c *Context) SetHeader(key string, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) String(code int64, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int64, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int64, data []byte) {
	c.Status(code)
	c.writer.Write(data)
}

func (c *Context) Param(name string) interface{} {
	for _, param := range c.Params {
		if param.Key == name {
			return param.Value
		}
	}
	return nil
}

func (c *Context) HTML(code int64, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.writer.Write([]byte(html))
}

func (c *Context) Clear() {

	c.Method = ""
	c.Path = ""
	c.Keys = nil
	c.writer = nil
	c.request = nil
	c.handlers = c.handlers[:0]
	c.Params = nil
	c.index = -1
	c.FullPath = ""
}
