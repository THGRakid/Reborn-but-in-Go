package submission

import (
	"Reborn-but-in-Go/submission/controller"
	"Reborn-but-in-Go/submission/dao"
	"Reborn-but-in-Go/submission/service"
	"github.com/gin-gonic/gin"
)

func InitSubmissionRouter(r *gin.Engine) {
	//创建数据访问层（DAO）的单例实例
	submissionDao := dao.NewSubmissionDaoInstance()

	//创建服务层（Service）的实例，传递数据访问层实例
	submissionService := service.NewSubmissionService(submissionDao)

	//创建表现层（controller）的实例，传递服务层实例
	submissionController := controller.NewSubmissionController(submissionService)

	//注册POST路由，登录用户选择视频上传，使用表现层的------函数。
	r.POST("/douyin/publish/action/", submissionController.CreateVideo)

	//注册GET路由，登录用户的视频发布列表，直接列出用户所有投稿过的视频。使用表现层的-----函数。
	r.GET("/douyin/publish/list/", submissionController.QueryVideoList)

}
