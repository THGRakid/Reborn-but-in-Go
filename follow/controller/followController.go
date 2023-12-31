package controller

import (
	"Reborn-but-in-Go/follow/dao"
	"Reborn-but-in-Go/follow/service"
	"Reborn-but-in-Go/middleware"
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

// RelationActionResp 关注和取消关注需要返回结构
type RelationActionResp struct {
	Response
}

// FollowingResp 关注列表相关部分
type FollowingResp struct {
	Response
	UserList []model.User `json:"user_list,omitempty"`
}

// FollowersResp 粉丝列表相关部分
type FollowersResp struct {
	Response
	UserList []model.User `json:"user_list,omitempty"`
}

// RelationAction 处理关注和取消关注请求。
func (f *FollowController) RelationAction(c *gin.Context) {
	// 验证Token
	middleware.AuthMiddleware()(c)
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证 token 获得的信息：", isAuthenticated)
	// token 验证失败
	if !isAuthenticated.(bool) {
		log.Println("token 验证失败")
		return
	}
	// token 验证通过，可以继续处理
	userIDInterface, _ := c.Get("user_id")
	userIdInt, err1 := userIDInterface.(int)
	if !err1 {
		// 类型转换失败
		fmt.Println("用户id格式错误")
		c.JSON(http.StatusInternalServerError, RelationActionResp{
			Response{
				StatusCode: 2,
				StatusMsg:  "用户id类型转换失败",
			},
		})
		return
	}
	userId := int64(userIdInt)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)

	// 传入参数格式有问题。
	if nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("参数有误")
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "传入参数格式错误",
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
		go dao.NewFollowDaoInstance().DeleteFollowRelation(userId, toUserId)
	}
	c.JSON(http.StatusOK, RelationActionResp{
		Response{
			StatusCode: 0,
			StatusMsg:  "关注/取关操作成功",
		},
	})
}

// GetFollowing 处理获取关注列表的请求
func (f *FollowController) GetFollowing(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	// 用户id解析出错
	if nil != err {
		fmt.Println("用户id格式错误")
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
			UserList: nil,
		})
		return
	}
	// 正常获取关注列表
	users, err := f.FollowService.GetFollowing(userId)
	if err != nil {
		// 获取关注列表时出错
		fmt.Println("获取关注列表时出错")
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表时出错",
			},
			UserList: nil,
		})
		return
	} else {
		// 获取关注列表成功
		log.Println("获取关注列表成功")
		c.JSON(http.StatusOK, FollowingResp{
			UserList: users,
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "获取关注列表成功",
			},
		})
		return
	}
}

// GetFollowers 处理获取粉丝列表的请求
func (f *FollowController) GetFollowers(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	// 用户id解析出错
	if nil != err {
		fmt.Println("用户id格式错误")
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
			UserList: nil,
		})
		return
	}
	// 正常获取粉丝列表
	users, err := f.FollowService.GetFollowers(userId)
	if err != nil {
		// 获取粉丝列表时出错
		fmt.Println("获取粉丝列表时出错")
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表时出错",
			},
			UserList: nil,
		})
		return
	} else {
		// 获取粉丝列表成功
		log.Println("获取粉丝列表成功")
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "获取粉丝列表成功",
			},
			UserList: users,
		})
		return
	}
}
