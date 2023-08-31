package service

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/video/dao"
	"Reborn-but-in-Go/video/model"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"path/filepath"
	"time"
)

type VideoService struct {
	VideoDao *dao.VideoDao
}

// NewVideoService 创建一个新的 videoService 实例
func NewVideoService(videoDao *dao.VideoDao) *VideoService {
	return &VideoService{
		VideoDao: videoDao,
	}
}

const videoNum = 5 //feed每次返回的视频数量
// FeedGet 获得视频列表
func FeedGet(lastTime int64) ([]model.Video, error) {
	//t := time.Now()
	//fmt.Println(t)
	if lastTime == 0 { //没有传入参数或者视屏已经刷完
		lastTime = time.Now().Unix()
	}
	strTime := fmt.Sprint(time.Unix(lastTime, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("查询的时间", strTime)
	var VideoList []model.Video
	VideoList = make([]model.Video, 0)
	err := config.DB.Model(&model.Video{}).Where("create_at < ?", strTime).Order("create_at desc").Limit(videoNum).Find(&VideoList).Error
	return VideoList, err
}

// GetCoverPath 传入视频源文件地址，截取视频的第num帧作为视频封面，返回封面路径，可直接存至数据库中
/*
ffmpeg需访问http://ffmpeg.org/download.html下载安装并且配置好环境变量，相应教程可在网上查找，我的参考
https://www.bilibili.com/video/BV1xf4y1Z7pV/?spm_id_from=333.337.search-card.all.click&vd_source=e63b0e97f9f156b04320cc032690c072
*/
func GetCoverPath(videoPath string, frameNum int) (coverPath string, err error) {
	//获取项目路径
	workPath, _ := os.Getwd()
	//获取视频完整文件名（有后缀）
	videoName := filepath.Base(videoPath)
	//去掉视频文件后缀用作封面名
	coverName := videoName[:len(videoName)-len(filepath.Ext(videoName))]
	//拼接封面存放位置
	coverPath = workPath + "/static/covers/" + coverName
	//开始制作视频封面
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败1：", err)
		return "", err
	}
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败2：", err)
		return "", err
	}
	//将图片保存至本地
	err = imaging.Save(img, coverPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败3：", err)
		return "", err
	}
	//返回相对路径
	coverPath = coverPath[len(workPath):] + ".png"
	return coverPath, nil
}
