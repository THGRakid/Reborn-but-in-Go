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
func (Message) TableName() string {
	return "messages"
}

type FriendUser struct {
	Message         string `json:"message"`                                                                                             // 和该好友的最新聊天消息
	MsgType         int64  `json:"msgType"`                                                                                             // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	Id              int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                                 // 用户id
	Name            string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                                              // 用户名称
	FollowCount     int64  `protobuf:"varint,3,opt,name=follow_count,json=followCount,proto3" json:"follow_count,omitempty"`            // 关注总数
	FollowerCount   int64  `protobuf:"varint,4,opt,name=follower_count,json=followerCount,proto3" json:"follower_count,omitempty"`      // 粉丝总数
	IsFollow        bool   `protobuf:"varint,5,opt,name=is_follow,json=isFollow,proto3" json:"is_follow,omitempty"`                     // true-已关注，false-未关注
	Avatar          string `protobuf:"bytes,6,opt,name=avatar,proto3" json:"avatar,omitempty"`                                          //用户头像
	BackgroundImage string `protobuf:"bytes,7,opt,name=background_image,json=backgroundImage,proto3" json:"background_image,omitempty"` //用户个人页顶部大图
	Signature       string `protobuf:"bytes,8,opt,name=signature,proto3" json:"signature,omitempty"`                                    //个人简介
	TotalFavorited  int64  `protobuf:"varint,9,opt,name=total_favorited,json=totalFavorited,proto3" json:"total_favorited,omitempty"`   //获赞数量
	WorkCount       int64  `protobuf:"varint,10,opt,name=work_count,json=workCount,proto3" json:"work_count,omitempty"`                 //作品数量
	FavoriteCount   int64  `protobuf:"varint,11,opt,name=favorite_count,json=favoriteCount,proto3" json:"favorite_count,omitempty"`     //点赞数量

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
