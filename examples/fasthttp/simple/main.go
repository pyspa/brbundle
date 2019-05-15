package main

import (
	"fmt"

	"go.pyspa.org/brbundle/brfasthttp"
	"github.com/valyala/fasthttp"

	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	fmt.Println("Listening at :8080")
	fmt.Println("You can access index.html at /index.html")
	fasthttp.ListenAndServe(":8080", brfasthttp.Mount())
}
