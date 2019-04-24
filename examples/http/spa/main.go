package main

import (
	"fmt"
	"net/http"

	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brhttp"
	_ "github.com/shibukawa/brbundle/examples"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("You can access index.html at any location")
	m.Handle("/",
		brhttp.Mount(brbundle.WebOption{
			SPAFallback: "index.html",
		}),
	)
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", m)
}