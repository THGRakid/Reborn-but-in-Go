package submission

import (
	"Reborn-but-in-Go/submission/controller"
	"Reborn-but-in-Go/submission/dao"
	"Reborn-but-in-Go/submission/service"
	"github.com/gin-gonic/gin"
)

func InitSubmissionRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public") // 创建一个默认的 Gin 路由引擎

	//创建数据访问层（DAO）的单例实例
	submissionDao := dao.NewVideoDaoInstance()

	//创建服务层（Service）的实例，传递数据访问层实例
	submissionService := service.NewVideoService(submissionDao)

	//创建表现层（controller）的实例，传递服务层实例
	submissionController := controller.NewVideoController(submissionService)

	//注册POST路由，登录用户选择视频上传，使用表现层的------函数。
	r.POST("/douyin/publish/action/", submissionController.CreateVideo)

	//注册GET路由，登录用户的视频发布列表，直接列出用户所有投稿过的视频。使用表现层的-----函数。
	r.GET("/douyin/publish/list/", submissionController.QueryVideoList)

}
