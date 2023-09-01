package controller

import (
	"Reborn-but-in-Go/favorite/service"
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
	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 64)
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
}

// GetFavoriteList 获取点赞列表;
func (f *FavoriteController) GetFavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	strCurId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	curId, _ := strconv.ParseInt(strCurId, 10, 64)
	Favorite := new(service.FavoriteService)
	videos, err := Favorite.GetFavoriteList(userId, curId)
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
}
