package main

import (
	"fmt"
	"net/http"

	"go.pyspa.org/brbundle/brhttp"
	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	fmt.Println("Listening at :8080")
	fmt.Println("You can access index.html at /index.html")
	http.ListenAndServe(":8080", brhttp.Mount())
}

