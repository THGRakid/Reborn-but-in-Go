package model

import (
	"time"
)

// Message 持久层结构块
type Message struct {
	Id       int64     `json:"id"`           //消息id
	UserId   int64     `json:"from_user_id"` //发送者用户id
	ToUserId int64     `json:"to_user_id"`   //接收者用户id
	Content  string    `json:"content"`      //消息内容
	CreateAt time.Time `json:"create_at"`    //发送时间
}

// TableName 修改表名映射
func (Message) TableName() string {
	return "messages"
}

type FriendUser struct {
	Message string `json:"message"` // 和该好友的最新聊天消息
	MsgType int64  `json:"msgType"` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}

// Response 响应状态码
type Response struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// ChatResponse  查看消息响应结构快
type ChatResponse struct {
	Response
	MessageList []*Message `json:"message_list"` // 消息列表
}

// ListResponse  查看消息列表响应结构快
type ListResponse struct {
	Response
	UserList []FriendUser `json:"user_list"` // 消息列表
}

// RelationActionRequest  发送消息POST请求接收结构块
type RelationActionRequest struct {
	Token      string // 用户鉴权token
	ToUserId   int64  // 对方用户id
	ActionType int32  // 1-发送消息
	Content    string // 消息内容
}
