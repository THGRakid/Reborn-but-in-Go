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

// Count
// 1、使用 video id 查询 Comment 数量
func (*CommentDao) Count(videoId int64) (int64, error) {
	log.Println("CommentDao-Count: running")
	var count int64

	// 数据库查询评论数量，注意要使用 model.Comment 结构体
	err := config.DB.Model(&model.Comment{}).Where(map[string]interface{}{"video_id": videoId, "status": 1}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed")
		return -1, err
	}
	log.Println("CommentDao-Count: return count success")
	return count, nil
}

// CommentIdList
// 2、根据视频id获取评论id列表
func (*CommentDao) CommentIdList(videoId int64) ([]int64, error) {
	var commentIdList []int64

	// 从数据库中查询评论id列表，注意要使用 model.Comment 结构体
	err := config.DB.Model(&model.Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentDao-CommentIdList: query comment id list failed")
		return nil, err
	}
	return commentIdList, nil
}

// InsertComment
// 3、发表评论
func (*CommentDao) InsertComment(comment model.Comment) (model.Comment, error) {
	log.Println("CommentDao-InsertComment: running")

	// 在评论结构体中设置默认值
	comment.Status = 1
	comment.CreateAt = time.Now()

	// 向数据库插入一条评论信息
	err := config.DB.Model(&model.Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: create comment failed")
		return model.Comment{}, err
	}
	log.Println("CommentDao-InsertComment: create comment success")
	return comment, nil
}

// DeleteComment
// 4、删除评论，传入评论id
func (*CommentDao) DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running")

	// 查询评论信息，检查是否存在且为有效评论
	var commentInfo model.Comment
	result := config.DB.Model(&model.Comment{}).Where(map[string]interface{}{"id": id, "status": 1}).First(&commentInfo)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-DeleteComment: comment does not exist")
		return errors.New("comment does not exist")
	}

	// 将评论状态更新为无效
	err := config.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", 0).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: delete comment failed")
		return err
	}
	log.Println("CommentDao-DeleteComment: delete comment success")
	return nil
}

// GetCommentList
// 5、根据视频id查询所属评论全部列表信息
func (*CommentDao) GetCommentList(videoId int64) ([]model.Comment, error) {
	log.Println("CommentDao-GetCommentList: running")

	// 查询评论信息列表，按时间倒序排列
	var commentList []model.Comment
	result := config.DB.Model(&model.Comment{}).Where(map[string]interface{}{"video_id": videoId, "status": 1}).
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
