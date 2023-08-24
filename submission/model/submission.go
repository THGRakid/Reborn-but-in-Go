package model

import (
	"time"
)

// Submission 持久层结构块
type Video struct {
	Id            int64     `gorm:"primaryKey"` //视频id
	UserId        int64     //作者id
	VideoPath     string    //视频地址
	CoverPath     string    //视频封面地址
	FavoriteCount int64     //点赞数
	CommentCount  int64     //评论数
	Title         string    //视频标题
	Time          time.Time //投稿时间
	Status        int32     //待审核0 发布成功1 审核失败2 下架删除3
}

// TableName 修改表名映射
func (Video) TableName() string {
	return "videos"
}

// 获取视频列表 结构块
type ListResponse struct {
	StatusCode int32    //状态码，0-成功，其他值-失败
	StatusMsg  string   //返回状态描述
	VideoList  []*Video //视频列表
}
