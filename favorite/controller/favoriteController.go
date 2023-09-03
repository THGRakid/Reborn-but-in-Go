package controller

import (
	"Reborn-but-in-Go/favorite/service"
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/user/model"
	vidDao "Reborn-but-in-Go/video/dao"
	vidMod "Reborn-but-in-Go/video/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FavoriteController struct {
	FavoriteService *service.FavoriteService
}

// NewFavoriteController 创建一个新的 FavoriteController 实例，并传递 FavoriteService
func NewFavoriteController(FavoriteService *service.FavoriteService) *FavoriteController {
	return &FavoriteController{
		FavoriteService: FavoriteService,
	}
}

type FavoriteResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavoriteListResponse struct {
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

// FavoriteAction 点赞或者取消赞操作;
func (f *FavoriteController) FavoriteAction(c *gin.Context) {
	middleware.AuthMiddleware()(c)
	//验证Token
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		// token 验证通过，可以继续处理

		//根据Token获取userId
		userIDInterface, _ := c.Get("user_id")
		userIdInt, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
			return
		}
		userId := int64(userIdInt)
		//从前端接受请求参数
		strVideoId := c.Query("video_id")
		videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
		strActionType := c.Query("action_type")
		actionType, _ := strconv.ParseInt(strActionType, 10, 64)
		//	调用服务层
		Favorite := new(service.FavoriteService)
		//获取点赞或者取消赞操作的错误信息
		err := Favorite.FavoriteAction(userId, videoId, int8(actionType))
		if err == nil && actionType == 1 {
			log.Printf("点赞成功")
			c.JSON(http.StatusOK, FavoriteResponse{
				StatusCode: 0,
				StatusMsg:  "favorite action success",
			})
		} else if err == nil && actionType == 2 {
			log.Printf("取消赞成功")
			c.JSON(http.StatusOK, FavoriteResponse{
				StatusCode: 0,
				StatusMsg:  "remove favorite action success",
			})
		} else {
			log.Printf("点赞失败：%v", err)
			c.JSON(http.StatusOK, FavoriteResponse{
				StatusCode: 1,
				StatusMsg:  "favorite action fail",
			})
		}
	} else {
		// token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}

// GetFavoriteList 获取点赞列表;
func (f *FavoriteController) GetFavoriteList(c *gin.Context) {
	middleware.AuthMiddleware()(c)
	//验证Token
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		// token 验证通过，可以继续处理

		//根据Token获取userId
		userIDInterface, _ := c.Get("user_id")
		userIdInt, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			c.JSON(http.StatusInternalServerError, GetFavoriteListResponse{
				StatusCode: 2,
				StatusMsg:  "Failed to convert user_id to int",
				VideoList:  nil, // 设置为空，因为返回数据类型要求是数组
			})
			return
		}
		userId := int64(userIdInt)
		//从前端接受请求参数
		Favorite := new(service.FavoriteService)
		videos, err := Favorite.GetFavoriteList(userId)
		videoModlist, err := Favorite.GetVideosByVideoIDs(videos)
		if err == nil {
			log.Printf("获取点赞列表成功")
			c.JSON(http.StatusOK, GetFavoriteListResponse{
				StatusCode: 0,
				StatusMsg:  "get favoriteList success",
				VideoList:  convertToVideoResponse(videoModlist), // 转换为新的 VideoResponse 结构
			})
		} else {
			log.Printf("获取点赞列表失败：%v", err)
			c.JSON(http.StatusOK, GetFavoriteListResponse{
				StatusCode: 1,
				StatusMsg:  "get favoriteList fail",
				VideoList:  nil, // 设置为空，因为返回数据类型要求是数组
			})
		}
	} else {
		// token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}

// 辅助函数，将 Video 转换为 VideoResponse
func convertToVideoResponse(videos []vidMod.Video) []VideoResponse {
	var videoResponses []VideoResponse
	Favorite := new(service.FavoriteService)
	for _, video := range videos {
		authorId, err := vidDao.GetVideoAuthor(video.Id)
		if err != nil {
			fmt.Printf("获取视频作者失败：%v\n", err)
		}
		user, err := Favorite.GetUserByID(authorId)
		if err != nil {
			fmt.Printf("获取用户信息失败：%v\n", err)
		}
		boolIsFavorite, err := service.IsFavorite(video.Id, authorId)
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
			IsFavorite:    boolIsFavorite,
			Title:         video.Title,
		})
	}
	return videoResponses
}
