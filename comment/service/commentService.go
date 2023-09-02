package service

import (
	"Reborn-but-in-Go/comment/dao"
	"Reborn-but-in-Go/comment/model"
	"errors"
)

// CommentService 评论服务层
type CommentService struct {
	CommentDao *dao.CommentDao
}

// NewCommentService 创建一个新的 CommentService 实例
func NewCommentService(commentDao *dao.CommentDao) *CommentService {
	return &CommentService{
		CommentDao: commentDao,
	}
}

// CreateComment 发表评论
func (s *CommentService) CreateComment(videoID, userID int64, content string) (model.Comment, error) {
	// 检查评论内容是否为空
	if content == "" {
		return model.Comment{}, errors.New("评论内容不能为空")
	}

	// 调用 CommentDao 的 InsertComment 方法插入评论
	newComment, err := s.CommentDao.InsertComment(videoID, userID, content)
	if err != nil {
		return model.Comment{}, err
	}

	return newComment, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentID, usrId int64) error {
	// 调用 CommentDao 的 DeleteComment 方法删除评论
	err := s.CommentDao.DeleteComment(commentID, usrId)
	return err
}

// GetCommentListByVideoId 根据视频ID查询评论列表
func (s *CommentService) GetCommentListByVideoId(videoID int64) ([]model.Comment, error) {
	// 调用 CommentDao 的 GetCommentList 方法获取评论列表
	commentList, err := s.CommentDao.GetCommentList(videoID)
	if err != nil {
		return nil, err
	}

	return commentList, nil
}
