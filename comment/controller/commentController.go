package controller

import (
	"Reborn-but-in-Go/comment/service"
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/user/model"
	userService "Reborn-but-in-Go/user/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CommentController 定义评论的控制器结构体
type CommentController struct {
	commentService *service.CommentService
	userService    *userService.UserService // 使用新包名
}

// NewCommentController 创建一个新的 CommentController 实例
func NewCommentController(commentService *service.CommentService, userService *userService.UserService) *CommentController {
	return &CommentController{
		commentService: commentService,
		userService:    userService,
	}
}

type CommentResponse struct {
	StatusCode int32   `json:"status_code"`
	StatusMsg  string  `json:"status_msg,omitempty"`
	Comment    Comment `json:"comment,omitempty"`
}

type Comment struct {
	ID         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// 其他用户信息字段
}

type GetCommentResponse struct {
	StatusCode  int32    `json:"status_code"`
	StatusMsg   string   `json:"status_msg,omitempty"`
	CommentList []string `json:"comment_list,omitempty"`
}

// HandleCommentAction 统一处理评论操作请求
func (controller *CommentController) HandleCommentAction(c *gin.Context) {
	// 验证Token
	middleware.AuthMiddleware()(c)
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {

		// token 验证通过，可以继续处理
		// 根据Token获取userId
		userIDInterface, _ := c.Get("user_id")
		userIdInt, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": "Failed to convert user_id to int"})
			return
		}
		userId := int64(userIdInt)

		// 从前端接受请求参数
		strVideoId := c.Query("video_id")
		videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
		strActionType := c.Query("action_type")
		actionType, _ := strconv.Atoi(strActionType)

		// 获取用户信息
		userResponse, err := controller.userService.GetUserByID(userId)
		if err != nil {
			// 处理获取用户信息失败的情况
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 5, "status_msg": "Failed to get user information"})
			return
		}

		switch actionType {
		case 1:
			commentText := c.DefaultQuery("comment_text", "") // 获取评论内容

			// 调用 CommentService 中的 CreateComment 方法发表评论
			newComment, err := controller.commentService.CreateComment(videoId, userId, commentText)
			if err != nil {
				c.JSON(http.StatusInternalServerError, CommentResponse{
					StatusCode: 1, // 设置适当的错误状态码
					StatusMsg:  "Failed to create comment",
				})
				return
			}

			// 返回新创建的评论
			c.JSON(http.StatusOK, CommentResponse{
				StatusCode: 0,
				Comment: Comment{
					ID: newComment.UserId,
					User: User{
						ID:   userResponse.User.Id,   // 使用评论的用户ID
						Name: userResponse.User.Name, // 使用从 userService 获取的用户名
					},
					Content:    newComment.Content,
					CreateDate: newComment.CreateAt.Format("01-02-2006"), // 格式化创建日期为 mm-dd 格式
				},
			})
		case 2:
			// 从请求参数中获取评论ID
			commentIDStr := c.Param("comment_id")
			commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
			if err != nil {
				// 处理评论ID无效的情况
				c.JSON(http.StatusBadRequest, gin.H{"status_code": 3, "status_msg": "Invalid comment ID"})
				return
			}

			// 调用 CommentService 中的 DeleteComment 方法删除评论
			err = controller.commentService.DeleteComment(commentID, userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status_code": 4, "status_msg": "Failed to delete comment"})
				return
			}

			// 返回成功删除评论的响应
			c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Comment deleted successfully"})
		}
	} else {
		// Token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}

}

// GetCommentList 修改 GetCommentList 函数，获取评论列表时关联用户信息
func (controller *CommentController) GetCommentList(c *gin.Context) {
	middleware.AuthMiddleware()(c)
	// 验证 Token
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证 token 获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		// Token 验证通过，可以继续处理

		// 根据 Token 获取 userId
		userIDInterface, _ := c.Get("user_id")
		_, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
			return
		}
		// 从前端接受请求参数
		strVideoId := c.Query("video_id")
		videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
		// 调用 CommentService 中的 GetCommentListByVideoId 方法获取评论列表
		commentList, err := controller.commentService.GetCommentListByVideoId(videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment list"})
			return
		}

		// 遍历评论列表，获取用户信息并合并到结果中
		var result []Comment
		for _, comment := range commentList {
			// 获取用户信息
			userResponse, err := controller.userService.GetUserByID(comment.UserId)
			if err != nil {
				// 处理获取用户信息失败的情况
				c.JSON(http.StatusInternalServerError, gin.H{"status_code": 5, "status_msg": "Failed to get user information"})
				return
			}

			// 创建包含用户信息的评论
			commentWithUser := Comment{
				ID: comment.UserId,
				User: User{
					ID:   userResponse.User.Id,
					Name: userResponse.User.Name,
					// 其他用户信息字段
				},
				Content:    comment.Content,
				CreateDate: comment.CreateAt.Format("01-02-2006"), // 格式化创建日期为 mm-dd-yyyy 格式
			}

			result = append(result, commentWithUser)
		}

		// 返回评论列表
		c.JSON(http.StatusOK, gin.H{"comment_list": result})
	} else {
		// Token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}
