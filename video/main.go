package main

import (
	"Reborn-but-in-Go/video/controller"
	"Reborn-but-in-Go/video/dao"
	"Reborn-but-in-Go/video/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	videoDao := dao.NewVideoDaoInstance()
	videoService := service.NewVideoService(videoDao)
	feedController := controller.NewFeedController(videoService)
	r.GET("/douyin/feed", feedController.Feed)
	if err := r.Run(":8080"); err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
