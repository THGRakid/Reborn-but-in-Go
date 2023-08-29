package comment

import (
	"Reborn-but-in-Go/comment/controller"
	"Reborn-but-in-Go/comment/dao"
	"Reborn-but-in-Go/comment/service"
	"github.com/gin-gonic/gin"
)

func InitCommentRouter(r *gin.Engine) {
	// 创建数据访问层（DAO）的单例实例
	commentDao := dao.NewCommentDaoInstance()

	// 创建服务层（Service）的实例，传递数据访问层实例
	commentService := service.NewCommentService(commentDao)

	// 创建表现层（Controller）的实例，传递服务层实例
	commentController := controller.NewCommentController(commentService)

	// 注册 GET 路由，查看视频的所有评论，按发布时间倒序，使用表现层中的 GetCommentList 函数
	r.GET("/douyin/comment/list/", commentController.GetCommentList)

	// 注册 POST 路由，登录用户对视频进行评论，使用表现层中的 HandleCommentAction 函数
	r.POST("/douyin/comment/action/", commentController.HandleCommentAction)

}
