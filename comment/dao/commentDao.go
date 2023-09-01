package dao

import (
	"Reborn-but-in-Go/comment/model"
	"Reborn-but-in-Go/config"
	"errors"
	"log"
	"sync"
	"time"
)

// CommentDao 评论的数据库访问对象（DAO）
type CommentDao struct {
}

// 用于保存实例
var commentDao *CommentDao

// 单例模式，只生成一个 commentDao 实例，提高性能
var commentOnce sync.Once

// NewCommentDaoInstance 用于获取 CommentDao 单例实例的函数
func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(func() {
		commentDao = &CommentDao{}
	})
	return commentDao
}

// InsertComment 插入评论
func (*CommentDao) InsertComment(videoID, userID int64, content string) (model.Comment, error) {
	log.Println("CommentDao-InsertComment: running")

	// 创建评论结构体实例，并设置默认值
	comment := model.Comment{
		VideoId:  videoID,
		UserId:   userID,
		Content:  content,
		Status:   1,
		CreateAt: time.Now(),
	}

	// 向数据库插入一条评论信息
	err := config.DB.Model(&model.Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: create comment failed")
		return model.Comment{}, err
	}
	log.Println("CommentDao-InsertComment: create comment success")
	return comment, nil
}

// DeleteComment 删除评论
func (*CommentDao) DeleteComment(commentID int64) error {
	log.Println("CommentDao-DeleteComment: running")

	// 查询评论信息，检查是否存在且为有效评论
	var commentInfo model.Comment
	result := config.DB.Model(&model.Comment{}).Where("id = ? AND status = ?", commentID, 1).First(&commentInfo)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-DeleteComment: comment does not exist")
		return errors.New("comment does not exist")
	}

	// 将评论状态更新为无效
	err := config.DB.Model(&model.Comment{}).Where("id = ?", commentID).Update("status", 0).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: delete comment failed")
		return err
	}
	log.Println("CommentDao-DeleteComment: delete comment success")
	return nil
}

// GetCommentList 获取评论列表
func (*CommentDao) GetCommentList(videoID int64) ([]model.Comment, error) {
	log.Println("CommentDao-GetCommentList: running")

	// 查询评论信息列表，按时间倒序排列
	var commentList []model.Comment
	result := config.DB.Model(&model.Comment{}).Where("video_id = ? AND status = ?", videoID, 1).
		Order("create_at desc").Find(&commentList)

	// 若没有评论信息，返回空列表而不是错误
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: no comments for this video")
		return nil, nil
	}

	// 若获取评论列表出错
	if result.Error != nil {
		log.Println("CommentDao-GetCommentList: get comment list failed")
		return commentList, result.Error
	}
	log.Println("CommentDao-GetCommentList: get comment list success")
	return commentList, nil
}
