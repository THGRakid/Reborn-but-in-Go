package controller

import (
	"Reborn-but-in-Go/message/model"
	"Reborn-but-in-Go/message/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// MessageController 表现层
type MessageController struct {
	MessageService *service.MessageService
	//RedisClient    *redis.Client
}

// 创建一个新的 MessageController 实例，并传递 MessageService
func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{
		MessageService: messageService,
	}
}

/*
func (c *MessageController) HandleWebSocketConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// 处理错误
		return
	}
	defer conn.Close()

	// 将 conn 加入到订阅列表，以便在有新消息时推送给客户端
	c.MessageService.AddSubscriber(conn)

	// 保持 WebSocket 连接，等待客户端发送或关闭连接
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// 处理错误或客户端关闭连接，从订阅列表中移除 conn
			c.MessageService.RemoveSubscriber(conn)
			break
		}

		// 在收到消息时进行处理，根据需要发送响应
		// 例如，可以根据消息内容进行一些操作，并将结果发送给客户端
		// ...

		// 示例：将消息内容返回给客户端
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			// 处理错误
		}
	}
}
*/
// QueryMessage 处理获取聊天消息的请求
func (c *MessageController) QueryMessage(ctx *gin.Context) {
	// 获取请求参数
	//token := ctx.Query("token")
	toUserIdString := ctx.Query("to_user_id")
	preMsgTime := ctx.Query("pre_msg_time")

	//将获取的string类型数据改成int64
	toUserId, _ := strconv.Atoi(toUserIdString)

	// 在这里进行用户鉴权，校验 token 等
	// ...还不会
	/*
		在进行认证或授权检查时，如果用户未通过验证或没有权限，可以中止请求的处理。

		if !isAuthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	*/

	// 调用服务层获取聊天消息记录
	chatResponse, err := c.MessageService.QueryMessage(int64(toUserId), preMsgTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat messages"})
		return
	}

	// 返回获取的聊天消息记录
	ctx.JSON(http.StatusOK, chatResponse)
}

// 处理消息操作的请求
func (c *MessageController) SendMessage(ctx *gin.Context) {
	// 绑定请求中的 JSON 数据到 RelationActionRequest 结构体
	var actionRequest model.RelationActionRequest
	if err := ctx.ShouldBindJSON(&actionRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 在这里进行用户鉴权，校验 token 等
	// ...还不会
	/*
		在进行认证或授权检查时，如果用户未通过验证或没有权限，可以中止请求的处理。

		if !isAuthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	*/

	// 检查 action_type，如果不是发送消息，返回错误
	if actionRequest.ActionType != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action type"})
		return
	}

	// 调用服务层发送消息
	err := c.MessageService.SendMessage(actionRequest.ToUserId, actionRequest.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// 返回操作成功的响应
	ctx.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Action completed successfully"})
}
