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
	CreateAt      time.Time //投稿时间
	Status        int32     //待审核0 发布成功1 审核失败2 下架删除3
}

// TableName 修改表名映射
func (Video) TableName() string {
	return "videos"
}

// User 持久层结构块
type User struct {
	Id              int64  `json:"id,omitempty"`   //用户id
	Name            string `json:"name,omitempty"` //用户名称（不重复）
	FollowCount     int64  `json:"follow_count"`   //关注人数
	FollowerCount   int64  `json:"follower_count"` //粉丝人数
	IsFollow        bool   `json:"is_follow"`      // true-未关注，false-已关注
	Password        string //密码
	Avatar          string `json:"avatar"`           //用户头像
	BackgroundImage string `json:"background_image"` //背景图像
	Signature       string `json:"signature"`        //个人简介
	TotalFavorited  int64  `json:"total_favorited"`  //获赞数
	WorkCount       int64  `json:"work_count"`       //作品数
	FavoriteCount   int64  `json:"favorite_count"`   //点赞数
}

// TableName 修改表名映射
func (User) TableName() string {
	return "users"
}
