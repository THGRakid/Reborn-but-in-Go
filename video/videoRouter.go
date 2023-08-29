package video

import (
	"Reborn-but-in-Go/video/controller"
	"Reborn-but-in-Go/video/dao"
	"Reborn-but-in-Go/video/service"
	"github.com/gin-gonic/gin"
)

func InitVideoRouter(r *gin.Engine) {
	videoDao := dao.NewVideoDaoInstance()

	videoService := service.NewVideoService(videoDao)

	feedController := controller.NewFeedController(videoService)

	r.GET("/douyin/feed/", feedController.Feed)

}
