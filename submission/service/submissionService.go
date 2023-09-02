package service

import (
	"Reborn-but-in-Go/submission/dao"
	"Reborn-but-in-Go/submission/model"
	"Reborn-but-in-Go/video/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// UploadFileToServer 传入要保存到服务器上的地址以及本地文件地址，保存文件到服务器并返回文件存储地址
func UploadFileToServer(key, path string) (finalPath string, err error) {
	accessKey := "LZPQ0Bx_xldazTQsnD1VYvnIe3aWSPhQsLGF9lML" //密钥
	secretKey := "uVj2nC2fn2KlkROZppBNPXOz6zrJpEV_J99ehZto" //密钥
	bucket := "zmxs"                                        //空间名称
	//生成上传凭证
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}

	err = formUploader.PutFile(context.Background(), &ret, upToken, key, path, &putExtra)

	return "http://s0apnlizm.hn-bkt.clouddn.com/" + key, err
}

// CreateVideo 保存视频到
func (s *SubmissionService) CreateVideo(userId int64, title string, data *multipart.FileHeader, ctx *gin.Context) error {
	//获取文件名
	videoName := filepath.Base(data.Filename)
	//将userid和视频文件名进行拼接得到最终视频文件名
	videoName = fmt.Sprintf("%d_%s", userId, videoName)

	videoPath := filepath.Join("../videos/", videoName)
	//将最终视频文件暂时保存至本地
	if err := ctx.SaveUploadedFile(data, videoPath); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to save video to host",
		})

		return err
	}

	//开始上传视频到服务器
	finalVideoPath, err := UploadFileToServer("videos/"+videoName, videoPath)
	if err != nil {

		ctx.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to upload video to server from submissionService",
		})
		fmt.Println("上传视频的err：", err)
		return err
	}

	//调用videoService的 GetCoverPath 函数，获取封面地址,并暂时保存到了本地
	coverPath, err := service.GetCoverPath(videoPath, 1)
	//失败则无法投稿
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to get covers from videoService",
		})
		return err
	}
	//开始上传封面到服务器
	coverName := strings.Replace(videoName, ".mp4", ".png", 1)
	finalCoverPath, err := UploadFileToServer("covers/"+coverName, coverPath)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to get covers from videoService",
		})
		fmt.Println("上传封面的err：", err)
		return err
	}

	//删除保存在本地中的文件
	err = os.Remove(videoPath)
	if err != nil {
		panic(err)
	}
	err = os.Remove(coverPath)
	if err != nil {
		panic(err)
	}
	//需要处理视频数据data，得到视频地址
	video := &model.Video{
		UserId:        userId,
		VideoPath:     finalVideoPath,
		CoverPath:     finalCoverPath,
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
