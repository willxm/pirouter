package main

import (
	"fmt"
	"time"

	"github.com/willxm/pirouter"
)

func main() {
	r := pirouter.NewRouter()
	r.Register("POST", "/user/login", greet)

	r.Register("GET", "/", hiRoot)
	r.Register("GET", "/user/info", hix)
	r.Register("GET", "/user/photo", hiy)
	r.Register("GET", "/user/:name", rawParam1)
	//r.Register("GET", "/user/:name2", rawParam)
	r.Register("GET", "/user/:name/2", rawParam2)
	r.Register("GET", "/user/:name/:age", rawParam3)
	r.Register("GET", "/user/:name/:age/3", rawParam4)

	r.Register("GET", "/user/mid", mid, hix)

	r.Register("GET", "/:name/info", hix)
	r.Register("GET", "/user", hiy)



	r.Run(":8080")
}

func greet(c *pirouter.Context) {

	fmt.Fprintf(c.Writer, "Hello! %s\n", time.Now())
}

func hiRoot(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi root! %s\n", time.Now())
}

func hix(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi x! %s\n", time.Now())
}

func hiy(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi y! %s\n", time.Now())
}

func rawParam1(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi rawParam1 %s\n", c.Req.URL)
}

func rawParam2(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi rawParam2 %s\n", c.Req.URL)
}


func rawParam3(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi rawParam3 %s\n", c.Req.URL)
}


func rawParam4(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Hi rawParam4 %s\n", c.Req.URL)
}


func mid(c *pirouter.Context) {
	fmt.Fprintf(c.Writer, "Mid Befor ! %s\n", time.Now())

	c.Next()

	fmt.Fprintf(c.Writer, "Mid After! %s\n", time.Now())
}
