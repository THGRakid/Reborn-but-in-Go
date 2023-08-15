package model

import (
	"time"
)

// Message 持久层结构块
type Message struct {
	Id        int64     `gorm:"primaryKey"` //消息id
	SendId    int64     //发送者用户id
	ReceiveId int64     //接收者用户id
	Content   string    //消息内容
	CreateAt  time.Time //发送时间
	Status    int8      //消息状态（存在为1，默认为0）
	DeletedAt time.Time //删除时间
}

// TableName 修改表名映射
func (Message) TableName() string {
	return "messages"
}
