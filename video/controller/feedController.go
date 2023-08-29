package controller

import (
	favoriteService "Reborn-but-in-Go/favorite/service"
	followService "Reborn-but-in-Go/follow/service"
	userDao "Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/video/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type feedController struct {
	VideoService    *service.VideoService
	FollowService   *followService.FollowService
	FavoriteService *favoriteService.FavoriteService
	UserDao         *userDao.UserDao
}

//——————————————中间件——————————————

var Key = []byte("Reborn_But_In_Go") //加密key
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

var userId int64 = 666
var userName = "GGBond"

// CreateToken 生成一个token
func CreateToken(userId int64, userName string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	nowTime := time.Now()                        //当前时间
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "zhoumo",          //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(Key)
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// ——————————————中间件——————————————

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

func (controller *feedController) Feed(c *gin.Context) {
	s, e := CreateToken(userId, userName)
	if e != nil {
		fmt.Printf("%s", e)
	}
	fmt.Println(s)
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
		var user, err = controller.UserDao.GetUserByID(v.UserId)
		var feedUser FeedUser
		if err == nil { //用户存在
			feedUser.Id = user.Id
			feedUser.Follower = user.Follower
			feedUser.Follow = user.Following
			feedUser.Name = user.Name
			//add
			feedUser.FavoritedCount = user.FavoritedCount
			feedUser.FavoriteCount = user.FavoriteCount
			feedUser.IsFollow = false
			if haveToken {
				// 查询是否关注
				tokenStruct, ok := CheckToken(strToken)               //中间件校验Token
				if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
					uid1 := tokenStruct.UserId                                      //用户id
					uid2 := v.UserId                                                //视频发布者id
					isFollow, _ := controller.FollowService.IsFollowing(uid1, uid2) //传入两个userId，检查是否关注
					if isFollow {
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
			tokenStruct, ok := CheckToken(strToken)
			if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
				uid := tokenStruct.UserId                                         //用户id
				vid := v.Id                                                       // 视频id
				isFavorite, _ := controller.FavoriteService.IsFavourite(vid, uid) //点赞，传入视频Id和userId，检查该用户是否点赞了此视频
				if isFavorite {                                                   //有点赞记录
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
