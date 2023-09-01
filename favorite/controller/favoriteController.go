package controller

import (
	"Reborn-but-in-Go/favorite/service"
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/user/model"
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
	StatusCode int32   `json:"status_code"`
	StatusMsg  string  `json:"status_msg,omitempty"`
	VideoList  []int64 `json:"video_list,omitempty"`
	//VideoList  []model.Video `json:"video_list,omitempty"`
	//错误情况，后续videoList使用类型为[]int64，而本需要[]model.Video ，需完成对接工作。未更改
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
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
			return
		}
		userId := int64(userIdInt)
		//从前端接受请求参数
		Favorite := new(service.FavoriteService)
		videos, err := Favorite.GetFavoriteList(userId)
		if err == nil {
			log.Printf("获取点赞列表成功")
			c.JSON(http.StatusOK, GetFavoriteListResponse{
				StatusCode: 0,
				StatusMsg:  "get favoriteList success",
				VideoList:  videos,
			})
		} else {
			log.Printf("获取点赞列表失败：%v", err)
			c.JSON(http.StatusOK, GetFavoriteListResponse{
				StatusCode: 1,
				StatusMsg:  "get favoriteList fail ",
			})
		}
	} else {
		// token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}
