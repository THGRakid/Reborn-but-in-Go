package config

import (
	"github.com/gin-gonic/gin"
)

func GetGinEngine() *gin.Engine {
	r := gin.Default()
	return r
}
