package model

import (
	"time"
)

// Message 持久层结构块
type Message struct {
	Id       int64     `gorm:"primaryKey"` //消息id
	UserId   int64     //发送者用户id
	ToUserId int64     //接收者用户id
	Content  string    //消息内容
	CreateAt time.Time //发送时间
	Status   int8      //消息状态（存在为1，默认为0）
}

// TableName 修改表名映射
func (Message) TableName() string {
	return "messages"
}

// ChatResponse  查看消息响应结构快
type ChatResponse struct {
	StatusCode  int32      // 状态码，0-成功，其他值-失败
	StatusMsg   string     // 返回状态描述
	MessageList []*Message // 消息列表
}

// RelationActionRequest  发送消息POST请求接收结构块
type RelationActionRequest struct {
	Token      string // 用户鉴权token
	ToUserId   int64  // 对方用户id
	ActionType int32  // 1-发送消息
	Content    string // 消息内容
}
