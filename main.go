package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"字节青训营/Reborn-but-in-Go/config"
)

func main() {
	r := config.GetGinEngine()

	prefix := "/api"

	r.GET(prefix+"/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from API!"})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
