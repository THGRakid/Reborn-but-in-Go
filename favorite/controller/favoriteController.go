package controller

//：表现层，对接前端的接口
import (
	"Reborn-but-in-Go/favorite/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// FavoriteController 表现层
type FavoriteController struct {
	FavoriteService *service.FavoriteService
}

// 创建一个新的 FavoriteController 实例，并传递 FavoriteService
func NewFavoriteController(FavoriteService *service.FavoriteService) *FavoriteController {
	return &FavoriteController{
		FavoriteService: FavoriteService,
	}
}

type FavoriteResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavouriteListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	VideoList  []service.Video `json:"video_list,omitempty"`
}

// FavoriteAction 点赞或者取消赞操作;
func (f *FavoriteController) FavoriteAction(c *gin.Context) {
	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 64)
	Favorite := new(service.FavoriteServiceImpl)
	//获取点赞或者取消赞操作的错误信息
	err := Favorite.FavouriteAction(userId, videoId, int32(actionType))
	if err == nil {
		log.Printf("方法Favorite.FavouriteAction(userid, videoId, int32(actiontype) 成功")
		c.JSON(http.StatusOK, FavoriteResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
		log.Printf("方法Favorite.FavouriteAction(userid, videoId, int32(actiontype) 失败：%v", err)
		c.JSON(http.StatusOK, FavoriteResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action fail",
		})
	}
}

// GetFavouriteList 获取点赞列表;
func (f *FavoriteController) GetFavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	strCurId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	curId, _ := strconv.ParseInt(strCurId, 10, 64)
	Favorite := GetVideo()
	videos, err := Favorite.GetFavouriteList(userId, curId)
	if err == nil {
		log.Printf("方法Favorite.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
			StatusMsg:  "get favouriteList success",
			VideoList:  videos,
		})
	} else {
		log.Printf("方法Favorite.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "get favouriteList fail ",
		})
	}
}
