package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/submission/model"
	"fmt"
	"sync"
)

/*
	？？？？构造方法？？？？
*/

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

// 2、视频列表。根据userID，查出video列表
func (*SubmissionDao) QueryVideoList(userId int64) ([]*model.Video, error) {
	var videos []*model.Video
	//根据user_id查询视频列表，按从大到小排序
	err := config.DB.Where("user_id = ?", userId).Order("time desc").Find(&videos).Error
	if err != nil {
		fmt.Println("Failed to get video list")
		return nil, err
	}
	return videos, nil
}
