package controller

import (
	favoriteService "Reborn-but-in-Go/favorite/service"
	followService "Reborn-but-in-Go/follow/service"
	"Reborn-but-in-Go/middleware"
	userDao "Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/video/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type feedController struct {
	VideoService    *service.VideoService
	FollowService   *followService.FollowService
	FavoriteService *favoriteService.FavoriteService
	UserDao         *userDao.UserDao
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
	Id            int64    `json:"id"`
	Author        FeedUser `json:"author"`
	VideoPath     string   `json:"play_url"`
	CoverPath     string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
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
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Follow          int64  `json:"follow_count"`
	Follower        int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}

func (controller *feedController) Feed(c *gin.Context) {

	middleware.AuthMiddleware()(c)
	//验证Token
	isAuthenticated, _ := c.Get("is_authenticated")
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
		//tmp.Author依靠用户信息接口查询
		var user, err = controller.UserDao.GetUserByID(v.UserId)
		var feedUser FeedUser
		if err == nil { //用户存在
			feedUser.Id = user.Id
			feedUser.Name = user.Name
			feedUser.Follower = user.FollowerCount
			feedUser.Follow = user.FollowCount
			feedUser.Avatar = user.Avatar
			feedUser.BackgroundImage = user.BackgroundImage
			feedUser.Signature = user.Signature
			feedUser.WorkCount = user.WorkCount
			feedUser.TotalFavorited = user.TotalFavorited
			feedUser.FavoriteCount = user.FavoriteCount
			feedUser.IsFollow = false

			//查询是否关注
			if isAuthenticated.(bool) {
				// token 验证通过，可以继续处理
				// 获取userId
				userIdString := c.Query("user_id")
				userId, _ := strconv.Atoi(userIdString)
				uid1 := int64(userId) //用户id
				uid2 := v.UserId      //视频发布者id
				fmt.Println("feed层获取到的视频发布者ID：", v.Id)
				isFollow, _ := controller.FollowService.IsFollowing(uid1, uid2) //传入两个userId，检查是否关注
				if isFollow {
					feedUser.IsFollow = true
				}
			}
		}
		tmp.Author = feedUser
		tmp.CommentCount = v.CommentCount
		tmp.CoverPath = v.CoverPath
		tmp.FavoriteCount = v.FavoriteCount
		tmp.IsFavorite = false
		//查询是否点赞过
		if isAuthenticated.(bool) {
			// token 验证通过，可以继续处理
			// 获取userId
			userIdString := c.Query("user_id")
			userId, _ := strconv.Atoi(userIdString)
			uid := int64(userId)                                             //用户id
			vid := v.Id                                                      // 视频id
			isFavorite, _ := controller.FavoriteService.IsFavorite(vid, uid) //点赞，传入视频Id和userId，检查该用户是否点赞了此视频
			if isFavorite {                                                  //有点赞记录
				tmp.IsFavorite = true
			}
		}

		tmp.Title = v.Title
		feedVideoList = append(feedVideoList, tmp)
		newTime = v.CreateAt.Unix()
	}
	if len(feedVideoList) > 0 {
		fmt.Println(feedVideoList)
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
