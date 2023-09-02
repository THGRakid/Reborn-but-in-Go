package comment

import (
	"Reborn-but-in-Go/comment/controller"
	"Reborn-but-in-Go/comment/dao"
	commentService "Reborn-but-in-Go/comment/service"
	userD "Reborn-but-in-Go/user/dao"
	userService "Reborn-but-in-Go/user/service"
	"github.com/gin-gonic/gin"
)

func InitCommentRouter(r *gin.Engine) {
	// 创建评论数据访问层（DAO）的单例实例
	commentDao := dao.NewCommentDaoInstance()

	// 创建评论服务层（Service）的实例，传递评论DAO实例
	commentServiceInstance := commentService.NewCommentService(commentDao)

	// 创建用户数据访问层（DAO）的单例实例
	userDao := userD.NewUserDaoInstance()

	// 创建用户服务层（Service）的实例，传递用户DAO实例
	userServiceInstance := userService.NewUserService(userDao)

	// 创建评论控制器的实例，传递评论服务层和用户服务层实例
	commentController := controller.NewCommentController(commentServiceInstance, userServiceInstance)

	// 注册路由
	r.GET("/douyin/comment/list/", commentController.GetCommentList)
	r.POST("/douyin/comment/action/", commentController.HandleCommentAction)
}
