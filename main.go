package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}

func main() {
	router := gin.Default()
	router.GET("/*url", handler)
	router.Run()
}
