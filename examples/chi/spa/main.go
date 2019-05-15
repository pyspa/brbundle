package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.pyspa.org/brbundle"
	"go.pyspa.org/brbundle/brchi"
	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	r := chi.NewRouter()
	r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("You can access index.html at any location")
	r.NotFound(brchi.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	}))
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", r)
}
