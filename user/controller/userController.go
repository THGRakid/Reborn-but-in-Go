package controller

import (
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/user/model"
	"Reborn-but-in-Go/user/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// userController 表现层
type UserController struct {
	UserService *service.UserService
}

// 创建一个新的 userController 实例，并传递 userService
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// CreateUser 注册新用户
func (c *UserController) CreateUser(ctx *gin.Context) {

	username := ctx.Query("username")
	password := ctx.Query("password")

	loginResponse, err := c.UserService.CreateUser(username, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "注册失败"})
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)

}

// UserLogin 用户登录
func (c *UserController) UserLogin(ctx *gin.Context) {

	username := ctx.Query("username")
	password := ctx.Query("password")

	userResponse, err := c.UserService.UserLogin(username, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "登录失败"})
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// GetUserByID 通过id返回用户信息
func (c *UserController) GetUserByID(ctx *gin.Context) {
	middleware.AuthMiddleware()(ctx)
	//验证Token
	isAuthenticated, _ := ctx.Get("is_authenticated")
	if isAuthenticated.(bool) {
		// token 验证通过，可以继续处理
		// 获取userId
		userIdString := ctx.Query("user_id")
		userId, _ := strconv.Atoi(userIdString)
		userResponse, err := c.UserService.GetUserByID(int64(userId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"StatusCode": 1, "StatusMsg": "获取Id失败"})
			return
		}

		ctx.JSON(http.StatusOK, userResponse)

	} else {
		// token 验证未通过，返回登录页面
		ctx.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token验证失败"},
		})
	}

}
