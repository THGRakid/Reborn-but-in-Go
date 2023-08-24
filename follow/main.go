package main

import (
	"Reborn-but-in-Go/follow/controller"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	followDao := dao.NewFollowDaoInstance()

	followService := service.NewFollowService(followDao)

	followController := controller.NewFollowController(followService)

	r.POST("/relation/action/", followController.RelationAction)

	r.GET("/relation/follow/list/", followController.GetFollowing)

	r.GET("/relation/follower/list", followController.GetFollowers)

	// 启动服务器并监听在 :8080 端口上
	if err := r.Run(":8080"); err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
