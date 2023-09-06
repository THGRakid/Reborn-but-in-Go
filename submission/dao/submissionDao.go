package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/submission/model"
	userMod "Reborn-but-in-Go/user/model"
	"errors"
	"log"
	"sync"
)

type SubmissionDao struct {
}

// 用于保存单例实例
var submissionDao *SubmissionDao

// 单例模式，只生成一个SubmissionDao实例，提高性能
var submissionOnce sync.Once

// NewSubmissionDaoInstance 用于获取SubmissionDao单例实例的函数
// 传递一个匿名函数（闭包），其作用是创建一个新的 MessageDao 实例并将其赋值给 messageDao 变量。
// 这个函数只会在第一次调用 Do 方法时执行
func NewSubmissionDaoInstance() *SubmissionDao {
	submissionOnce.Do(
		func() {
			submissionDao = &SubmissionDao{}
		})
	return submissionDao
}

// 1、视频发布。创建一个新的video。
func (*SubmissionDao) CreateVideo(video *model.Video) error {
	//将video内数据导入数据库
	result := config.DB.Create(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 2、视频列表。根据userID，查出video发布列表
// QueryVideoList 根据userId查询所投稿的全部videoId
func (*SubmissionDao) QueryVideoList(userId int64) ([]int64, error) {
	var VideoIdList []int64
	err := config.DB.Model(model.Video{}).Where(map[string]interface{}{"user_id": userId}).
		Pluck("id", &VideoIdList).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("there are no PublishVideoId")
			return VideoIdList, nil
		} else {
			log.Println(err.Error())
			return VideoIdList, errors.New("Failed to get video ID list")
		}
	}
	return VideoIdList, nil
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userId int64) (userMod.User, error) {
	user := userMod.User{}
	err := config.DB.Model(userMod.User{}).Where(map[string]interface{}{"id": userId}).
		First(&user).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return userMod.User{}, nil
		} else {
			log.Println(err.Error())
			return user, errors.New("get userInfo failed")
		}
	}
	return user, nil
}
