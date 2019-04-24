package main

import (
	"fmt"
	_ "github.com/shibukawa/brbundle/examples"

	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle/brecho"
)

func main() {
	e := echo.New()
	fmt.Println("You can access index.html at /index.html")
	e.GET("/*", brecho.Mount())
	e.Logger.Fatal(e.Start(":1323"))
}
