package controller

import (
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

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	userResponse, err := c.UserService.CreateUser(username, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// UserLogin 用户登录
func (c *UserController) UserLogin(ctx *gin.Context) {

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	userResponse, err := c.UserService.UserLogin(username, password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// GetUserByID 通过id返回用户信息
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userIdString := ctx.Query("user_id")

	//将获取的string类型数据改成int64
	userId, _ := strconv.Atoi(userIdString)

	idResponse, err := c.UserService.GetUserByID(int64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	ctx.JSON(http.StatusOK, idResponse)
}
