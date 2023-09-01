package controller

import (
	"Reborn-but-in-Go/comment/service"
	"Reborn-but-in-Go/middleware"
	"Reborn-but-in-Go/user/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CommentController 定义评论的控制器结构体
type CommentController struct {
	commentService *service.CommentService
}

// NewCommentController 创建一个新的 CommentController 实例
func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

type CommentResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

type GetCommentResponse struct {
	StatusCode  int32    `json:"status_code"`
	StatusMsg   string   `json:"status_msg,omitempty"`
	CommentList []string `json:"comment_list,omitempty"`
}

// HandleCommentAction 统一处理评论操作请求
func (controller *CommentController) HandleCommentAction(c *gin.Context) {
	// 验证Token
	isAuthenticated, exists := c.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if !exists {
		// 处理未设置 "is_authenticated" 键的情况，例如返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 2,
			"status_msg":  "Missing authentication information",
		})
		return
	}
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
				Comment:    newComment.Content, // 获取评论内容
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
			err = controller.commentService.DeleteComment(commentID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status_code": 4, "status_msg": "Failed to delete comment"})
				return
			}

			// 返回成功删除评论的响应
			c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "Comment deleted successfully"})

		}
	} else {
		// token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}

}

// GetCommentList 处理获取评论列表的请求
func (controller *CommentController) GetCommentList(c *gin.Context) {
	middleware.AuthMiddleware()(c)
	//验证Token
	isAuthenticated, _ := c.Get("is_authenticated")
	fmt.Println("验证token获得的信息：", isAuthenticated)
	if isAuthenticated.(bool) {
		// token 验证通过，可以继续处理

		//根据Token获取userId
		userIDInterface, _ := c.Get("user_id")
		_, ok := userIDInterface.(int)
		if !ok {
			// 类型转换失败
			// 这里你可以处理转换失败的情况，例如返回错误信息
			fmt.Println("Error: Failed to convert user_id to int")
			c.JSON(http.StatusInternalServerError, gin.H{"status_code": 2, "status_msg": " Failed to convert user_id to int"})
			return
		}
		//从前端接受请求参数
		strVideoId := c.Query("video_id")
		videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
		// 调用 CommentService 中的 GetCommentListByVideoId 方法获取评论列表
		commentList, err := controller.commentService.GetCommentListByVideoId(videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment list"})
			return
		}

		// 返回评论列表
		c.JSON(http.StatusOK, gin.H{"comment_list": commentList})
	} else {
		// token 验证未通过，返回登录页面
		c.JSON(http.StatusOK, &model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Token authentication failed"},
		})
	}
}
