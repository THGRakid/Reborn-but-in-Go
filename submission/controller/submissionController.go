package controller

import (
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/submission/service"
	"Reborn-but-in-Go/user/model"
	vidDao "Reborn-but-in-Go/video/dao"
	vidMod "Reborn-but-in-Go/video/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetPublishListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	VideoList  []VideoResponse `json:"video_list,omitempty"`
}

// VideoResponse 结构体表示 video_list 中的单个视频信息
type VideoResponse struct {
	Id            int64        `json:"id"`
	Author        UserResponse `json:"author"`
	PlayURL       string       `json:"play_url"`
	CoverURL      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

// 修改原有的 UserResponse 结构体
type UserResponse struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}

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
}

// 2、处理 获取视频列表 的请求
func (c *SubmissionController) QueryVideoList(ctx *gin.Context) {
	middleware.AuthMiddleware()(ctx)
	//验证token
	isAuthenticated, _ := ctx.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		//token验证通过，继续处理
		//根据Token获取userId
		userIDInterface, _ := ctx.Get("user_id")
		userIdInt, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			ctx.JSON(http.StatusInternalServerError, GetPublishListResponse{
				StatusCode: 2,
				StatusMsg:  "Failed to convert user_id to int",
				VideoList:  nil, //空
			})
			return
		}
		userId := int64(userIdInt)
		videos, err := c.SubmissionService.QueryVideoList(userId)
		videoModelList, err := c.SubmissionService.GetVideosByVideoIDs(videos)
		if err != nil {
			fmt.Println("controller: Failed to get video list. ERR:#{err}")
			ctx.JSON(http.StatusOK, GetPublishListResponse{
				StatusCode: 1,
				StatusMsg:  "get publishList fail",
				VideoList:  nil, //空
			})
		} else {
			fmt.Println("Success: get video list")
			ctx.JSON(http.StatusOK, GetPublishListResponse{
				StatusCode: 0,
				StatusMsg:  "Success: get video list",
				VideoList:  convertToVideoResponse(videoModelList),
			})
		}

	} else {
		//token验证失败
		ctx.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}

// 辅助函数，将 Video 转换为 VideoResponse
func convertToVideoResponse(videos []vidMod.Video) []VideoResponse {
	var videoResponses []VideoResponse
	Favorite := new(service.SubmissionService)

	for _, video := range videos {
		authorId, err := vidDao.GetVideoAuthor(video.Id)
		if err != nil {
			fmt.Printf("获取视频作者失败：%v\n", err)
		}
		user, err := Favorite.GetUserByID(authorId)
		if err != nil {
			fmt.Printf("获取用户信息失败：%v\n", err)
		}
		//boolIsFavorite, err := service.IsFavorite(video.Id, authorId)
		if err != nil {
			fmt.Printf("获取点赞状态失败：%v\n", err)
		}
		videoResponses = append(videoResponses, VideoResponse{
			Id: video.Id,
			Author: UserResponse{
				Id:              authorId,
				Name:            user.Name,
				FollowCount:     user.FollowCount,
				FollowerCount:   user.FollowerCount,
				IsFollow:        user.IsFollow,
				Avatar:          user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature:       user.Signature,
				TotalFavorited:  user.TotalFavorited,
				WorkCount:       user.WorkCount,
				FavoriteCount:   user.FavoriteCount,
			},
			PlayURL:       video.VideoPath,
			CoverURL:      video.CoverPath,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			//IsFavorite:    boolIsFavorite,
			Title: video.Title,
		})
	}
	return videoResponses
}
