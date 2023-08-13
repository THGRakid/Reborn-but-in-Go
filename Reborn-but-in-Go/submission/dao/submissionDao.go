package dao

import (
	"errors"
	"log"
	"time"
)

type Video struct {
	Id            int64     //视频id
	UserId        int64     //作者id
	VideoPath     string    //视频地址
	CoverPath     string    //视频封面地址
	FavoriteCount int64     //点赞数
	CommentCount  int64     //评论数
	Title         string    //视频标题
	Time          time.Time //投稿时间
	Status        int32     //待审核0 发布成功1 审核失败2 下架删除3
}

func (Video) TableName() string {
	return "videos"
}

// 1、 视频发布
func InsertVideo(video Video) (Video, error) {
	log.Println("SubmissionDao-InsertVideo: running")
	//往数据库中插入一条 视频 信息
	err := db.Model(Video{}).Create(&video)
	if err != nil {
		log.Println("SubmissionDao-InsertVideo: return create video failed") //打印错误信息
		return Video{}, errors.New("create video failed")
	}
	log.Println("SubmissionDao-InsertVideo: return success")
	return video, nil
}

// 2、 视频投稿列表
func GetSubmittedVideoList(userId int64) ([]Video, error) {
	// 函数已运行
	log.Println("SubmissionDao-GetSubmittedVideoList: running")
	var videoList []Video
	//在数据库中查询
	result := db.Model(Video{}).Where(map[string]interface{}{"user_id": userId}).
		Order("time desc").Find(&videoList)
	//若该用户没有投稿视频，返回空列表，不报错
	if result.RowAffected == 0 {

		return nil, nil
	}
	//获取投稿视频出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		return videoList, errors.New("get submitted video list failed")
	}
	return videoList, nil
}
