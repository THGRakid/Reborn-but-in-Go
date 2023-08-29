package user

import (
	"Reborn-but-in-Go/user/controller"
	"Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/user/service"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public") // 创建一个默认的 Gin 路由引擎

	//创建数据访问层（DAO）的单例实例
	userDao := dao.NewUserDaoInstance()

	//创建服务层（Service）的实例，传递数据访问层实例
	userService := service.NewUserService(userDao)

	//创建表现层（controller）的实例，传递服务层实例
	userController := controller.NewUserController(userService)

	//注册POST路由，新用户注册时提供用户名（需要保证唯一），密码。创建成功后返回用户 id 和权限token。
	r.POST("/douyin/user/register/", userController.CreateUser)

	//注册POST路由，通过用户名和密码进行登录，登录成功后返回用户 id 和权限 token.
	r.POST("/douyin/user/login/", userController.UserLogin)

	//注册GET路由，通过用户id，返回用户信息。
	r.GET("/douyin/user/", userController.GetUserByID)

}
