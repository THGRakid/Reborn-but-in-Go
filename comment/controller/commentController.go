package controller

import (
	"Reborn-but-in-Go/comment/model"
	"Reborn-but-in-Go/comment/service"
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

type DouyinCommentActionRequest struct {
	Token       string `json:"token" binding:"required"`
	VideoID     int64  `json:"video_id" binding:"required"`
	ActionType  int32  `json:"action_type" binding:"required"`
	CommentText string `json:"comment_text,omitempty"`
	CommentID   int64  `json:"comment_id,omitempty"`
}

// CreateComment 处理发表评论的请求
func (controller *CommentController) CreateComment(c *gin.Context) {
	var comment model.Comment

	// 解析请求中的 JSON 数据到 comment 结构体中
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// TODO: 验证用户身份和权限，检查 token 是否有效

	// 调用 CommentService 中的 CreateComment 方法发表评论
	newComment, err := controller.commentService.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 返回新创建的评论
	c.JSON(http.StatusOK, newComment)
}

// DeleteComment 处理删除评论的请求
func (controller *CommentController) DeleteComment(c *gin.Context) {
	commentId, err := strconv.ParseInt(c.Param("commentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// TODO: 验证用户身份和权限，检查 token 是否有效

	// 调用 CommentService 中的 DeleteComment 方法删除评论
	err = controller.commentService.DeleteComment(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// 返回删除成功的消息
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// HandleCommentAction 统一处理评论操作请求
func (controller *CommentController) HandleCommentAction(c *gin.Context) {
	var request DouyinCommentActionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if request.ActionType == 1 {
		controller.CreateComment(c)
	} else if request.ActionType == 2 {
		controller.DeleteComment(c)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action type"})
	}
}

// GetCommentList 处理获取评论列表的请求
func (controller *CommentController) GetCommentList(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// 调用 CommentService 中的 GetCommentListByVideoId 方法获取评论列表
	commentList, err := controller.commentService.GetCommentListByVideoId(videoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment list"})
		return
	}

	// 返回评论列表
	c.JSON(http.StatusOK, gin.H{"comment_list": commentList})
}
