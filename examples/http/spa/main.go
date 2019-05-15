package main

import (
	"fmt"
	"net/http"

	"go.pyspa.org/brbundle"
	"go.pyspa.org/brbundle/brhttp"
	_ "go.pyspa.org/brbundle/examples"
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