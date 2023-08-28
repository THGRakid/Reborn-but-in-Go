package controller

import (
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"Reborn-but-in-Go/user/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FollowController struct {
	FollowService *service.FollowService
}

func NewFollowController(followService *service.FollowService) *FollowController {
	return &FollowController{
		FollowService: followService,
	}
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

// FollowingResp 关注列表相关部分，User传入预留
type FollowingResp struct {
	Response
	UserList []model.User `json:"user_list,omitempty"`
}

// FollowersResp 粉丝列表相关部分，User传入预留
type FollowersResp struct {
	Response
	UserList []model.User `json:"user_list,omitempty"`
}

// RelationAction 处理关注和取消关注请求。
func (f *FollowController) RelationAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)
	// 传入参数格式有问题。
	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("失败")
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
		})
		return
	}
	// 正常处理
	switch {
	// 关注
	case 1 == actionType:
		go dao.NewFollowDaoInstance().InsertFollowRelation(userId, toUserId)
	// 取关
	case 2 == actionType:
		go dao.NewFollowDaoInstance().InsertFollowRelation(userId, toUserId)
	}
	log.Println("关注、取关成功。")
	c.JSON(http.StatusOK, RelationActionResp{
		Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

// GetFollowing 处理获取关注列表请求。
func (f *FollowController) GetFollowing(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		fmt.Printf("用户id格式错误")
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			UserList: nil,
		})
		return
	}
	// 正常获取关注列表
	users, err := f.FollowService.GetFollowing(userId)
	// 获取关注列表时出错。
	if err != nil {
		fmt.Printf("获取关注列表时出错")
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表时出错。",
			},
			UserList: nil,
		})
		return
	} else {
		// 成功获取到关注列表。
		log.Println("获取关注列表成功。")
		c.JSON(http.StatusOK, FollowingResp{
			UserList: users,
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "OK",
			},
		})
		return
	}
}

// GetFollowers 处理获取关注列表请求
func (f *FollowController) GetFollowers(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			UserList: nil,
		})
		return
	}
	// 正常获取粉丝列表
	users, err := f.FollowService.GetFollowers(userId)
	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	// 成功获取到粉丝列表。
	//log.Println("获取粉丝列表成功。")
	c.JSON(http.StatusOK, FollowersResp{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: users,
	})
}
