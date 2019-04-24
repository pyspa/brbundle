package main

import (
	"fmt"
	"github.com/shibukawa/brbundle"
	_ "github.com/shibukawa/brbundle/examples"

	"net/http"

	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle/brecho"
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

