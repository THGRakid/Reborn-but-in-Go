package service

import (
	"Reborn-but-in-Go/submission/dao"
	"Reborn-but-in-Go/submission/model"
	"fmt"
	"time"
)

// SubmissionService 服务层

type SubmissionService struct {
	SubmissionDao *dao.SubmissionDao
}

// NewSubmissionService 创建一个新的SubmissionService实例
func NewSubmissionService(submissionDao *dao.SubmissionDao) *SubmissionService {
	return &SubmissionService{
		SubmissionDao: submissionDao,
	}
}

// 假定视频地址和封面地址，如何获取呢？
var VideoPath string = ""
var CoverPath string = ""

// 1、投稿视频 ？？？data怎么处理？？？
func (s *SubmissionService) CreateVideo(userId int64, data []byte, title string) error {
	video := &model.Video{
		UserId:        userId,
		VideoPath:     VideoPath,
		CoverPath:     CoverPath,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		CreateAt:      time.Now(),
		Status:        1,
	}

	//调用 DAO 的CreateVideo 方法来保存消息到数据库
	err := s.SubmissionDao.CreateVideo(video)
	if err != nil {
		return err
	}
	return nil
}

// 2、获取视频列表 根据用户ID获取视频列表
func (s *SubmissionService) QueryVideoList(userId int64) (*model.ListResponse, error) {
	//调用 VideoDao 的 QueryVideoList 方法获取视频状态码，0-成功，其他值-失败列表
	videoList, err := s.SubmissionDao.QueryVideoList(userId)

	if err != nil {
		fmt.Println("Service:Failed to get video list")
	}

	//构建ListResponse对象，将查询到的消息记录填充进去
	listResponse := &model.ListResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		VideoList:  videoList,
	}
	return listResponse, nil
}
