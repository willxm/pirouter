package main

import (
	"errors"
	"fmt"

	"github.com/willxm/pirouter"
)

func main() {
	r := pirouter.NewRouter()
	r.Register("GET", "/user/photo", hiy)
	go func() {
		r.Run(":9091")
	}()
	r.Use("GET", "/user", new(exampleMiddleware))
	<-make(chan int)
}

type exampleMiddleware struct {
}

func (*exampleMiddleware) HandleRequest(c *pirouter.Context) error {
	c.Writer.Write([]byte("sd"))
	fmt.Println("testMiddlewareHere")
	return errors.New("fake error")
}
