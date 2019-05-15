package main

import (
	"fmt"
	_ "go.pyspa.org/brbundle/examples"

	"github.com/labstack/echo"
	"go.pyspa.org/brbundle/brecho"
)

func main() {
	e := echo.New()
	fmt.Println("You can access index.html at /index.html")
	e.GET("/*", brecho.Mount())
	e.Logger.Fatal(e.Start(":1323"))
}
