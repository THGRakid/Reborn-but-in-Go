package main

import (
	"net/http"
	"Reborn-but-in-Go/config"
	"github.com/gin-gonic/gin"
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
