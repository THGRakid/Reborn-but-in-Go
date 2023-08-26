package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/video/model"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once //单例模式，只生成一个VideoDao实例，提高性能

// NewVideoDaoInstance 返回VideoDao实例
func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

//// GetVideoList 按投稿时间倒序返回视频列表
//func (*VideoDao) GetVideoList(num int) ([]model.Video, error) {
//	var videos []model.Video
//	result := config.DB.Model(&model.Video{}).Order("time desc").Limit(num).Find(&videos)
//	if result.Error != nil {
//		fmt.Println("按投稿时间降序获取Video列表失败")
//		return nil, result.Error
//	}
//	return videos, nil
//}

// CreateVideo 创建一个新的Video，返回Video实例
func (*VideoDao) CreateVideo(video *model.Video) (*model.Video, error) {
	result := config.DB.Create(&video)
	if result.Error != nil {
		fmt.Println("创建Video失败")
		return nil, result.Error
	}
	return video, nil
}

// GetVideoListByUserId 根据userId获取用户发布的视频列表
func (*VideoDao) GetVideoListByUserId(userId int64) ([]model.Video, error) {
	var videos []model.Video
	if err := config.DB.Where("user_id = ?", userId).Find(&videos).Error; err != nil {
		fmt.Println("按用户id查询Video列表失败")
		return nil, err
	}
	return videos, nil
}

// GetVideoById 根据videoId，获取视频实体
func (*VideoDao) GetVideoById(videoId int64) (*model.Video, error) {
	video := model.Video{Id: videoId}

	result := config.DB.Where("video_id = ?", videoId).First(&video)
	err := result.Error
	if err != nil {
		fmt.Println("根据videoId查询Video实体失败")
		return nil, err
	}
	return &video, err
}

// AddCommentCount 根据videoId将评论条数+1
func AddCommentCount(videoId int64) error {
	if err := config.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		fmt.Println("评论条数+1失败")
		return err
	}
	return nil
}

// ReduceCommentCount 根据videoId将评论条数-1
func ReduceCommentCount(videoId int64) error {
	if err := config.DB.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		fmt.Println("评论条数-1失败")
		return err
	}
	return nil
}

// GetVideoAuthor 根据videoId得到视频作者
func GetVideoAuthor(videoId int64) (int64, error) {
	var video model.Video
	if err := config.DB.Model(&model.Video{}).Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.Id, err
	}
	return video.UserId, nil
}
