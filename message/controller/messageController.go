package controller

import (
	"Reborn-but-in-Go/message/model"
	"Reborn-but-in-Go/message/service"
	"Reborn-but-in-Go/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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

// QueryFriend处理获取聊天消息的请求
func (c *MessageController) QueryFriend(ctx *gin.Context) {
	// 获取请求参数
	// 验证Token
	middleware.AuthMiddleware()(ctx)
	isAuthenticated, _ := ctx.Get("is_authenticated")
	fmt.Println("验证 token 获得的信息：", isAuthenticated)
	// token 验证失败
	if !isAuthenticated.(bool) {
		log.Println("token 验证失败")
		return
	}

	userIdString := ctx.Query("user_id")

	//将获取的string类型数据改成int64
	userId, _ := strconv.Atoi(userIdString)

	// 调用服务层获取聊天消息记录
	listResponse, err := c.MessageService.QueryList(int64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "获取消息列表失败"})
		return
	}

	// 返回获取的聊天消息记录

	ctx.JSON(http.StatusOK, listResponse)
}

// QueryMessage 处理获取聊天消息的请求
func (c *MessageController) QueryMessage(ctx *gin.Context) {
	// 获取请求参数
	// 验证Token
	middleware.AuthMiddleware()(ctx)
	isAuthenticated, _ := ctx.Get("is_authenticated")
	fmt.Println("验证 token 获得的信息：", isAuthenticated)
	// token 验证失败
	if !isAuthenticated.(bool) {
		log.Println("token 验证失败")
		return
	}
	// token 验证通过，可以继续处理
	userIDInterface, _ := ctx.Get("user_id")
	userIdInt, err1 := userIDInterface.(int)
	if !err1 {
		// 类型转换失败
		fmt.Println("用户id格式错误")
		// token 验证未通过，返回登录页面
		ctx.JSON(http.StatusOK, &model.ChatResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token验证失败"},
		})
	}
	userId := int64(userIdInt)
	toUserIdString := ctx.Query("to_user_id")
	preMsgTime := ctx.Query("pre_msg_time")

	//将获取的string类型数据改成int64
	toUserId, _ := strconv.Atoi(toUserIdString)

	// 调用服务层获取聊天消息记录
	chatResponse, err := c.MessageService.QueryMessage(userId, int64(toUserId), preMsgTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "获取消息失败"})
		return
	}

	// 返回获取的聊天消息记录
	ctx.JSON(http.StatusOK, chatResponse)
}

// SendMessage 处理消息操作的请求
func (c *MessageController) SendMessage(ctx *gin.Context) {
	// 绑定请求中的 JSON 数据到 RelationActionRequest 结构体
	actionTypeString := ctx.Query("action_type")
	actionType, _ := strconv.Atoi(actionTypeString)
	// 检查 action_type，如果不是发送消息，返回错误
	if actionType != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"StatusCode": 1, "StatusMsg": "这不是发送消息"})
		return
	}

	// 验证Token
	middleware.AuthMiddleware()(ctx)
	isAuthenticated, _ := ctx.Get("is_authenticated")
	fmt.Println("验证 token 获得的信息：", isAuthenticated)
	// token 验证失败
	if !isAuthenticated.(bool) {
		log.Println("token 验证失败")
		return
	}
	// token 验证通过，可以继续处理
	userIDInterface, _ := ctx.Get("user_id")
	userIdInt, err1 := userIDInterface.(int)
	if !err1 {
		// 类型转换失败
		fmt.Println("用户id格式错误")
		// token 验证未通过，返回登录页面
		ctx.JSON(http.StatusOK, &model.ChatResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token验证失败"},
		})
	}
	userId := int64(userIdInt)

	toUserIdString := ctx.Query("to_user_id")
	toUserId, _ := strconv.Atoi(toUserIdString)
	content := ctx.Query("content")
	fmt.Println(content, "这是内容")
	// 调用服务层发送消息
	err := c.MessageService.SendMessage(userId, int64(toUserId), content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "发送消息失败"})
		return
	}

	// 返回操作成功的响应
	ctx.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "发送成功"})
}
