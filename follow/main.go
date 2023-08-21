package main

import (
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"Reborn-but-in-Go/message/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	followDao := dao.NewFollowDaoInstance()

	followService := service.NewFollowService(followDao)

	followController := controller.NewFollowController(followService)

	// 注册 GET 路由，处理获取聊天消息的请求，使用表现层中的 QueryMessageList 函数
	r.GET("/douyin/message/chat", followController.Queryfollow)

	// 注册 POST 路由，处理发送消息操作的请求，使用表现层中的 SendMessage 函数
	r.POST("/douyin/message/action", followController.Sendfollow)

	// 启动服务器并监听在 :8080 端口上
	if err := r.Run(":8080"); err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
