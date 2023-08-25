package service

import (
	"Reborn-but-in-Go/comment/dao"
	"Reborn-but-in-Go/comment/model"
	"time"
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

// CountCommentsByVideoId 根据视频ID查询评论数量
func (s *CommentService) CountCommentsByVideoId(videoId int64) (int64, error) {
	return s.CommentDao.Count(videoId)
}

// GetCommentIdListByVideoId 根据视频ID获取评论ID列表
func (s *CommentService) GetCommentIdListByVideoId(videoId int64) ([]int64, error) {
	return s.CommentDao.CommentIdList(videoId)
}

// CreateComment 发表评论
func (s *CommentService) CreateComment(comment model.Comment) (model.Comment, error) {
	// 设置评论的默认属性
	comment.Time = time.Now()
	comment.Status = 1 // 设置默认状态为有效

	// 调用 CommentDao 的 InsertComment 方法插入评论
	newComment, err := s.CommentDao.InsertComment(comment)
	if err != nil {
		return model.Comment{}, err
	}

	return newComment, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentId int64) error {
	// 调用 CommentDao 的 DeleteComment 方法删除评论
	err := s.CommentDao.DeleteComment(commentId)
	return err
}

// GetCommentListByVideoId 根据视频ID查询评论列表
func (s *CommentService) GetCommentListByVideoId(videoId int64) ([]model.Comment, error) {
	// 调用 CommentDao 的 GetCommentList 方法获取评论列表
	commentList, err := s.CommentDao.GetCommentList(videoId)
	if err != nil {
		return nil, err
	}

	return commentList, nil
}
