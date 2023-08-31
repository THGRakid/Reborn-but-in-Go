package model

import "time"

// Follow 关注模块
type Follow struct {
	Id         int64     // 自增ID
	UserID     int64     // 关注人ID
	FollowerID int64     // 被关注人ID
	CreateAt   time.Time // 关注发起时间
	Status     int8      // 关注状态（1：关注，默认0）
}

// TableName 修改表名映射
func (Follow) TableName() string {
	return "follows"
}
