package dao

import (
	"Reborn-but-in-Go/config"
	userModel "Reborn-but-in-Go/user/model"
	"Reborn-but-in-Go/video/model"
	"fmt"
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

	result := config.DB.Where("id = ?", videoId).First(&video)
	err := result.Error
	if err != nil {
		fmt.Println("根据videoId查询Video实体失败")
		return nil, err
	}
	return &video, err
}

// GetVideoAuthor 根据videoId得到视频作者
func GetVideoAuthor(videoId int64) (int64, error) {
	var video model.Video
	if err := config.DB.Model(&model.Video{}).Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.Id, err
	}
	return video.UserId, nil
}

// GetPublishCount 根据userid获得发布作品数
func (*VideoDao) GetPublishCount(userId int64) (int64, error) {
	var count int64
	fmt.Println("获取用户作品数咯~")
	err := config.DB.Model(&userModel.User{}).Pluck("work_count", &count).Where("id=?", userId).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}
