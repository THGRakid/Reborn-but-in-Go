package message

import (
	"Reborn-but-in-Go/message/controller"
	"Reborn-but-in-Go/message/dao"
	"Reborn-but-in-Go/message/service"
	"github.com/gin-gonic/gin"
)

func InitMessageRouter(r *gin.Engine) {
	// 创建数据访问层（DAO）的单例实例
	messageDao := dao.NewMessageDaoInstance()

	// 创建服务层（Service）的实例，传递数据访问层实例
	messageService := service.NewMessageService(messageDao)

	// 创建表现层（Controller）的实例，传递服务层实例
	messageController := controller.NewMessageController(messageService)

	// 注册 GET 路由，处理获取聊天消息列表的请求，使用表现层中的 QueryMessage 函数
	r.GET("/douyin/relation/friend/list/", messageController.QueryFriend)

	// 注册 GET 路由，处理获取聊天消息的请求，使用表现层中的 QueryMessage 函数
	r.GET("/douyin/message/chat/", messageController.QueryMessage)

	// 注册 GET路由，处理发送消息操作的请求，使用表现层中的 SendMessage 函数
	r.GET("/douyin/message/action/", messageController.SendMessage)

}
