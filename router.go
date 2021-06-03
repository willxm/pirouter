package pirouter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var httpMethod = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
	"PATCH":   true,
	"ANY":     true,
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Router struct {
	middlewares []Middleware
	roots       map[string]*Tree
}

func NewRouter() Router {
	rm := make(map[string]*Tree, 0)
	return Router{roots: rm}
}

func (r *Router) Use(method string, path string, middlewares ...Middleware) error {
	method = strings.ToUpper(method)
	if v, ok := httpMethod[method]; !ok || !v {
		panic(fmt.Sprintf("method:%s not support\n", method))
	}
	if _, ok := r.roots[method]; !ok {
		return errors.New("invalid path:" + path)
	}
	node, _ := r.roots[method].Find(path)
	if node == nil {
		return errors.New("invalid path:" + path)
	}
	node.middlewares = append(node.middlewares, middlewares...)
	return nil
}

func (r *Router) Register(method string, path string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	if v, ok := httpMethod[method]; !ok || !v {
		panic(fmt.Sprintf("method:%s not support\n", method))
	}
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = NewTree()
	}
	r.roots[method].Add(path, handler)
}

func (r *Router) getRoute(method string, path string) (*Node, []Middleware) {
	if _, ok := r.roots[method]; !ok {
		return nil, nil
	}
	return r.roots[method].Find(path)
}

func (r *Router) handle(c *Context) {
	node, middlewares := r.getRoute(c.Method, c.Path)
	if node == nil {
		c.NotFound()
		return
	}
	var path = c.Path
	if node.path != "/" {
		path = TrimPathPrefix(path)
	}
	if path == node.path {
		for _, m := range middlewares {
			if err := m.HandleRequest(c); err != nil {
				return
			}
		}
		node.handle.(HandlerFunc)(c.Writer, c.Req)
	} else {
		c.NotFound()
	}
}

// Run defines the method to start a http server
func (r *Router) Run(addr string) (err error) {
	log.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handle(newContext(w, req))
}
