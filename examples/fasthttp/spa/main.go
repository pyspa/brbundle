package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"go.pyspa.org/brbundle"
	"go.pyspa.org/brbundle/brfasthttp"
	"github.com/valyala/fasthttp"

	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	r := fasthttprouter.New()
	r.GET("/api/status", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Hello, World!")
	})
	fmt.Println("You can access index.html at any location")
	r.NotFound = brfasthttp.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	})
	fmt.Println("Listening at :8080")
	fasthttp.ListenAndServe(":8080", r.Handler)
}
