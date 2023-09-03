package follow

import (
	"Reborn-but-in-Go/follow/controller"
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"github.com/gin-gonic/gin"
)

func InitFollowRouter(r *gin.Engine) {
	// 创建数据访问层（DAO）的单例实例
	followDao := dao.NewFollowDaoInstance()

	// 创建服务层（Service）的实例，传递数据访问层实例
	followService := service.NewFollowService(followDao)

	// 创建表现层（controller）的实例，传递服务层实例
	followController := controller.NewFollowController(followService)

	// 注册POST路由，处理关注和取消关注的请求
	r.POST("/douyin/relation/action/", followController.RelationAction)

	// 注册GET路由，处理获取关注列表的请求
	r.GET("/douyin/relation/follow/list/", followController.GetFollowing)

	// 注册GET路由，处理获取粉丝列表的请求
	r.GET("/douyin/relation/follower/list", followController.GetFollowers)

}
