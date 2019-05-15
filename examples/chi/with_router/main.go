package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.pyspa.org/brbundle/brchi"
	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	r := chi.NewRouter()
	fmt.Println("You can access index.html at /public/index.html")
	r.Get("/public/*", brchi.Mount())
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", r)
}
