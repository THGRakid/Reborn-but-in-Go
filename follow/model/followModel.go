package model

import "time"

// Follow 关注模块
type Follow struct {
	Id         int64     // 自增ID
	UserId     int64     // 关注人ID
	FollowerId int64     // 被关注人ID
	CreateAt   time.Time // 关注发起时间
	Status     int8      // 关注状态 默认0，1为取消关注
}

// TableName 修改表名映射
func (Follow) TableName() string {
	return "follows"
}
