package controller

import (
	"Reborn-but-in-Go/submission/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// VideoController 表现层
type VideoController struct {
	VideoService *service.VideoService
}

// 创建一个新的 VideoController 实例，并传递 VideoService
func NewVideoController(videoService *service.VideoService) *VideoController {
	return &VideoController{
		VideoService: videoService,
	}
}

// 1、处理 视频投稿 的请求
func (c *VideoController) CreateVideo(ctx *gin.Context) {
	//获取请求参数
	// token := ctx.Query("token")
	data := []byte(ctx.Query("data"))
	title := ctx.Query("title")

	//进行用户鉴权

	//调用服务层
	err := c.VideoService.CreateVideo(data, title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit video"})
		return
	}
	//操作成功
	ctx.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Submit video successfully"})
}

// 2、处理 获取视频列表 的请求
func (c *VideoController) QueryVideoList(ctx *gin.Context) {
	//获取请求参数
	// token := ctx.Query("token")
	userIdString := ctx.Query("user_id")
	// 获取的string转换成int64
	userId, _ := strconv.ParseInt(userIdString, 10, 64)

	//用户token鉴权

	//调用服务层获取视频列表
	listResponse, err := c.VideoService.QueryVideoList(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video list"})
		return
	}

	// 正常，返回获取的视频列表
	ctx.JSON(http.StatusOK, listResponse)
}
