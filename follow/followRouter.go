package main

import (
	"Reborn-but-in-Go/follow/controller"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"github.com/gin-gonic/gin"
)

func initFollowRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	followDao := dao.NewFollowDaoInstance()

	followService := service.NewFollowService(followDao)

	followController := controller.NewFollowController(followService)

	r.POST("/relation/action/", followController.RelationAction)

	r.GET("/relation/follow/list/", followController.GetFollowing)

	r.GET("/relation/follower/list", followController.GetFollowers)

}
