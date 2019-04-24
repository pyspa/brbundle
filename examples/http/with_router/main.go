package main

import (
	"fmt"
	"net/http"

	"github.com/shibukawa/brbundle/brhttp"
	_ "github.com/shibukawa/brbundle/examples"
)

func main() {
	m := http.NewServeMux()
	fmt.Println("You can access index.html at /public/index.html")
	m.Handle("/public/", http.StripPrefix("/public", brhttp.Mount()))
	m.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", m)
}
