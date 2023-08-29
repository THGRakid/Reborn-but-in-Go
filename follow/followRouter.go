package follow

import (
	"Reborn-but-in-Go/follow/controller"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"github.com/gin-gonic/gin"
)

func InitFollowRouter(r *gin.Engine) {
	followDao := dao.NewFollowDaoInstance()

	followService := service.NewFollowService(followDao)

	followController := controller.NewFollowController(followService)

	r.POST("/relation/action/", followController.RelationAction)

	r.GET("/relation/follow/list/", followController.GetFollowing)

	r.GET("/relation/follower/list", followController.GetFollowers)

}
