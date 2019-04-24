package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brchi"
	_ "github.com/shibukawa/brbundle/examples"
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
