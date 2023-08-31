package service

import (
	"Reborn-but-in-Go/submission/dao"
	"Reborn-but-in-Go/submission/model"
	"Reborn-but-in-Go/video/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
//var VideoPath string = ""

// 1、投稿视频 ？？？data怎么处理？？？
func (s *SubmissionService) CreateVideo(userId int64, title string, data *multipart.FileHeader, ctx *gin.Context) error {
	//获取文件名
	videoName := filepath.Base(data.Filename)
	//将userid和视频文件名进行拼接得到最终视频文件名
	videoName = fmt.Sprintf("%d_%s", userId, videoName)
	//将最终视频文件保存至本地
	workPath, _ := os.Getwd()
	videoPath := workPath + "/static/videos/" + videoName
	if err := ctx.SaveUploadedFile(data, videoPath); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status_code": 1, "status_msg": "Failed to save video to host"})
		return err
	}
	//调用videoService的 GetCoverPath 函数，获取封面地址
	CoverPath, err := service.GetCoverPath(videoPath, 1)
	//失败则无法投稿
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to get covers from videoService",
		})
		return err
	}
	//需要处理视频数据data，得到视频地址
	video := &model.Video{
		UserId:        userId,
		VideoPath:     videoPath,
		CoverPath:     CoverPath,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		CreateAt:      time.Now(),
		Status:        1,
	}

	//调用 DAO 的CreateVideo 方法来保存消息到数据库
	err2 := s.SubmissionDao.CreateVideo(video)
	if err2 != nil {
		return err2
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
