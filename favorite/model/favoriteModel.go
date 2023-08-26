package model

import "time"

// Favorite 持久层结构块
type Favorite struct {
	UserId  int64
	VideoId int64
	Time    time.Time //点赞时间
	Status  int8      //点赞状态（点赞为1，取消赞为2）
}

// TableName 修改表名映射
func (Favorite) TableName() string {
	return "favorite"
}
