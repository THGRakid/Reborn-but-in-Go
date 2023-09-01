package model

import (
	"time"
)

// Comment 持久层结构块
type Comment struct {
	UserId   int64     //用户id
	VideoId  int64     //视频id
	Content  string    //评论内容
	CreateAt time.Time //评论发布时间
	Status   int8      //评论状态（存在为1，默认为0）
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comment"
}
