package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/key-1212/preflop", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		println(c.Request.Header.Get("User-Agent"))
	})
	r.Run()
}
