package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brgin"
	_ "github.com/shibukawa/brbundle/examples"
)

func main() {
	r := gin.Default()
	fmt.Println("You can access index.html at /public/index.html")
	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.NoRoute(brgin.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	}))
	fmt.Println("Listening at :8080")
	r.Run(":8080")
}
