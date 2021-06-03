package pirouter

//Middleware which implement Middleware could return non-nil error to finish request handler, such as non-permission;
type Middleware interface {
	HandleRequest(c *Context) error
}
