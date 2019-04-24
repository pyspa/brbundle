package main

import (
	"fmt"
	"net/http"

	"github.com/shibukawa/brbundle/brhttp"
	_ "github.com/shibukawa/brbundle/examples"
)

func main() {
	fmt.Println("Listening at :8080")
	fmt.Println("You can access index.html at /index.html")
	http.ListenAndServe(":8080", brhttp.Mount())
}

