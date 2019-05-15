package main

import (
	"fmt"
	"go.pyspa.org/brbundle"
	_ "go.pyspa.org/brbundle/examples"

	"net/http"

	"github.com/labstack/echo"
	"go.pyspa.org/brbundle/brecho"
)

func main() {
	e := echo.New()
	e.GET("/api/status", func (c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	fmt.Println("You can access index.html at any location")
	echo.NotFoundHandler = brecho.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	})
	e.Logger.Fatal(e.Start(":1323"))
}

