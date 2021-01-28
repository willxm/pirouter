package pirouter

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Router struct {
	roots map[string]*Tree
}

func NewRouter() Router {
	rm := make(map[string]*Tree, 0)
	return Router{roots: rm}
}

func (r *Router) Register(method string, path string, hander HandlerFunc) {
	//TODO: check method
	// see https://tools.ietf.org/html/rfc7231#section-4
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = NewTree()
	}
	r.roots[method].Add(path, hander)
}
