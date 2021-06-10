package pirouter

import (
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
}

type HandlerFunc func(c *Context)

type Router struct {
	roots map[string]*Tree
}

func NewRouter() Router {
	rm := make(map[string]*Tree, 0)
	return Router{roots: rm}
}

func (r *Router) Register(method string, path string, handlers ...HandlerFunc) {
	method = strings.ToUpper(method)
	if v, ok := httpMethod[method]; !ok || !v {
		panic(fmt.Sprintf("method:%s not support\n", method))
	}
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = NewTree()
	}
	r.roots[method].Add(path, handlers...)
}

func (r *Router) getRoute(method string, path string) []*Node {
	if _, ok := r.roots[method]; !ok {
		return nil
	}
	return r.roots[method].Find(path)
}

func (r *Router) handle(c *Context) {
	ns := r.getRoute(c.Method, c.Path)
	if len(ns) > 0 {
		c.handlers = ns[0].handle
		c.Next()
	} else {
		c.NotFound()
	}
}

// Run defines the method to start a http server
func (r *Router) Run(addr string) (err error) {
	//todo debug switch
	{
		for k, v := range r.roots {
			fmt.Println(k)
			v.String()
		}
	}
	log.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handle(newContext(w, req))
}
