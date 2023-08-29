package message

import (
	"Reborn-but-in-Go/message/controller"
	"Reborn-but-in-Go/message/dao"
	"Reborn-but-in-Go/message/service"
	"github.com/gin-gonic/gin"
)

func InitMessageRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public") // 创建一个默认的 Gin 路由引擎

	// 创建数据访问层（DAO）的单例实例
	messageDao := dao.NewMessageDaoInstance()

	// 创建服务层（Service）的实例，传递数据访问层实例
	messageService := service.NewMessageService(messageDao)

	// 创建表现层（Controller）的实例，传递服务层实例
	messageController := controller.NewMessageController(messageService)

	// 注册 GET 路由，处理获取聊天消息的请求，使用表现层中的 QueryMessage 函数
	r.GET("/douyin/message/chat", messageController.QueryMessage)
	//r.GET("/douyin/message/chat/ws", messageController.HandleWebSocketConnection)

	// 注册 POST 路由，处理发送消息操作的请求，使用表现层中的 SendMessage 函数
	r.POST("/douyin/message/action", messageController.SendMessage)

}
