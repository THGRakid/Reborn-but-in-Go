package controller

import (
	userDao "Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/video/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type feedController struct {
	VideoService *service.VideoService
}

// NewFeedController 创建一个新的 FeedController 实例
func NewFeedController(videoService *service.VideoService) *feedController {
	return &feedController{
		VideoService: videoService,
	}
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FeedVideo struct {
	Id            int64    `json:"id,omitempty"`
	Author        FeedUser `json:"author,omitempty"`
	VideoPath     string   `json:"video_path,omitempty"`
	CoverPath     string   `json:"cover_path,omitempty"`
	FavoriteCount int64    `json:"favorite_count,omitempty"`
	CommentCount  int64    `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
	Title         string   `json:"title,omitempty"`
}
type FeedResponse struct {
	Response
	VideoList []FeedVideo `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}
type FeedNoVideoResponse struct {
	Response
	NextTime int64 `json:"next_time"`
}
type FeedUser struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Follow         int64  `json:"follow_count,omitempty"`
	Follower       int64  `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
	FavoritedCount int64  `json:"favorited_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}

func (*feedController) Feed(c *gin.Context) {
	strToken := c.Query("token")
	var haveToken bool
	if strToken == "" {
		haveToken = false
	} else {
		haveToken = true
	}
	var strLastTime = c.Query("latest_time")
	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
	if err != nil {
		lastTime = 0
	}

	var feedVideoList []FeedVideo
	feedVideoList = make([]FeedVideo, 0)
	videoList, _ := service.FeedGet(lastTime)
	var newTime int64 = 0 //返回的视频的最久的一个的时间
	for _, v := range videoList {
		var tmp FeedVideo
		tmp.Id = v.Id
		tmp.VideoPath = v.VideoPath
		//tmp.Author = //依靠用户信息接口查询
		var user, err = userDao.GetUserByID(v.UserId)
		var feedUser FeedUser
		if err == nil { //用户存在
			feedUser.Id = user.Id
			feedUser.Follower = user.Follower
			feedUser.Follow = user.Following
			feedUser.Name = user.Name
			//add
			feedUser.FavoritedCount = user.Favorited_count
			feedUser.FavoriteCount = user.Favorite_count
			feedUser.IsFollow = false
			if haveToken {
				// 查询是否关注
				tokenStruct, ok := middleware.CheckToken(strToken)
				if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
					var uid1 = tokenStruct.UserId //用户id
					var uid2 = v.UserId           //视频发布者id
					if service.IsFollowing(uid1, uid2) {
						feedUser.IsFollow = true
					}
				}
			}
		}
		tmp.Author = feedUser
		tmp.CommentCount = v.CommentCount
		tmp.CoverPath = v.CoverPath
		tmp.FavoriteCount = v.FavoriteCount
		tmp.IsFavorite = false
		if haveToken {
			//查询是否点赞过
			tokenStruct, ok := middleware.CheckToken(strToken)
			if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
				var uid = tokenStruct.Id             //用户id
				var vid = v.Id                       // 视频id
				if service.CheckFavorite(uid, vid) { //有点赞记录
					tmp.IsFavorite = true
				}
			}
		}
		tmp.Title = v.Title
		feedVideoList = append(feedVideoList, tmp)
		newTime = v.Time.Unix()
	}
	if len(feedVideoList) > 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0}, //成功
			VideoList: feedVideoList,
			NextTime:  newTime,
		})
	} else {
		c.JSON(http.StatusOK, FeedNoVideoResponse{
			Response: Response{StatusCode: 0}, //成功
			NextTime: 0,                       //重新循环
		})
	}

}
