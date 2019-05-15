package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.pyspa.org/brbundle/brgin"

	_ "go.pyspa.org/brbundle/examples"
)

func main() {
	r := gin.Default()
	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	fmt.Println("You can access index.html at /static/index.html")
	r.GET("/static/*filepath", brgin.Mount())
	r.Run(":8080")
}
