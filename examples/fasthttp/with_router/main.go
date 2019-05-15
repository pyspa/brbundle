package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"go.pyspa.org/brbundle/brfasthttp"
	"github.com/valyala/fasthttp"

	_ "go.pyspa.org/brbundle/examples"
)

func main() {

	r := fasthttprouter.New()
	r.GET("/api/status", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Hello, World!")
	})
	fmt.Println("You can access index.html at /static/index.html")
	r.GET("/static/*filepath", brfasthttp.Mount())

	fmt.Println("Listening at :8080")
	fasthttp.ListenAndServe(":8080", r.Handler)
}
