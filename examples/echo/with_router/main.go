package main

import (
	"fmt"
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
	g := e.Group("/assets")
	fmt.Println("You can access index.html at /assets/index.html")
	g.GET("/*", brecho.Mount())
	e.Logger.Fatal(e.Start(":1323"))
}
