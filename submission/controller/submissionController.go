package controller

import (
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/submission/service"
	"Reborn-but-in-Go/user/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// VideoController 表现层
type SubmissionController struct {
	SubmissionService *service.SubmissionService
}

// 创建一个新的 VideoController 实例，并传递 VideoService
func NewSubmissionController(videoService *service.SubmissionService) *SubmissionController {
	return &SubmissionController{
		SubmissionService: videoService,
	}
}

// 1、处理 视频投稿 的请求
func (c *SubmissionController) Publish(ctx *gin.Context) {
	middleware.AuthMiddleware()(ctx)
	//验证Token
	isAuthenticated, _ := ctx.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		// token 验证通过，可以继续处理

		//根据Token获取userId
		userIDInterface, _ := ctx.Get("user_id")
		userIdInt, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
			return
		}
		userId := int64(userIdInt)
		//从前端接收请求参数
		data, err := ctx.FormFile("data")
		title := ctx.PostForm("title")
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"status_code": 1, "status_msg": err.Error()})
			return
		}
		//	进行用户鉴权

		//	调用服务层
		if err := c.SubmissionService.CreateVideo(userId, title, data, ctx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": 1, "status_msg": "Failed to submit video"})
			return
		}
		//	操作成功
		ctx.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Submit video successfully"})
	} else {
		// token 验证未通过，返回登录页面
		ctx.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}

	/*
		//验证Token
		isAuthenticated, _ := ctx.Get("is_authenticated")
		if isAuthenticated.(bool) {
			// token 验证通过，可以继续处理

			//根据Token获取userId
			userIDInterface, _ := ctx.Get("user_id")
			userId, ok := userIDInterface.(int64)
			if !ok {
				// 类型转换失败
				// 这里你可以处理转换失败的情况，例如返回错误信息
				fmt.Println("Error: Failed to convert user_id to int")
				ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
				return
			}
			//获取请求参数
			data := []byte(ctx.Query("data"))
			title := ctx.Query("title")

			//进行用户鉴权

			//调用服务层
			err := c.SubmissionService.CreateVideo(userId, data, title)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": 1, "status_msg": "Failed to submit video"})
				return
			}
			//操作成功
			ctx.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Submit video successfully"})
		} else {
			// token 验证未通过，返回登录页面
			ctx.JSON(http.StatusOK, &model.UserResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
			})
		}
	*/
}

// 2、处理 获取视频列表 的请求
func (c *SubmissionController) QueryVideoList(ctx *gin.Context) {
	//获取请求参数
	// token := ctx.Query("token")
	userIdString := ctx.Query("user_id")
	// 获取的string转换成int64
	userId, _ := strconv.ParseInt(userIdString, 10, 64)

	//用户token鉴权

	//调用服务层获取视频列表
	listResponse, err := c.SubmissionService.QueryVideoList(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video list"})
		return
	}

	// 正常，返回获取的视频列表
	ctx.JSON(http.StatusOK, listResponse)
}
