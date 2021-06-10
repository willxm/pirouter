package pirouter

import "net/http"

type Context struct {
	Writer   http.ResponseWriter
	Req      *http.Request
	Path     string
	Method   string
	Params   map[string]string //TODO: params encode
	handlers []HandlerFunc
	index    int8
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) NotFound() {
	c.Writer.WriteHeader(404)
	c.Writer.Write([]byte("NOT FOUND"))
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// TODO abort , last
