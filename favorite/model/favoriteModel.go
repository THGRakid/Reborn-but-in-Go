package model

import "time"

// Favorite 持久层结构块
type Favorite struct {
	UserId    int64
	VideoId   int64
	CreatedAt time.Time
}

// TableName 修改表名映射
func (Favorite) TableName() string {
	return "favorite"
}
