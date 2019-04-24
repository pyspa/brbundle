package main

import (
	"fmt"

	"github.com/shibukawa/brbundle/brfasthttp"
	"github.com/valyala/fasthttp"

	_ "github.com/shibukawa/brbundle/examples"
)

func main() {
	fmt.Println("Listening at :8080")
	fmt.Println("You can access index.html at /index.html")
	fasthttp.ListenAndServe(":8080", brfasthttp.Mount())
}
