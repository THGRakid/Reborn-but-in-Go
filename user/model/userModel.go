package model

import (
	"time"
)

// User 持久层结构块
type User struct {
	Id             int64     `gorm:"primaryKey"` //用户id
	Name           string    `gorm:"unique"`     //用户名称（不重复）
	Follower       int64     //粉丝人数
	Following      int64     //关注人数
	Password       string    //密码
	Avatar         string    //用户头像
	Background     string    //背景图像
	Introduce      string    //个人简介
	FavoritedCount int64     //获赞数
	WorkCount      int64     //作品数
	FavoriteCount  int64     //点赞数
	CreateAt       time.Time //用户创建时间
	Status         int8      `gorm:"default:0"` //用户状态（在线为1，不在线为0）
}

type UserForUsed struct {
	Id              int64  //用户id
	Name            string //用户名称（不重复）
	FollowCount     int64  //关注人数
	FollowerCount   int64  //粉丝人数
	IsFollow        bool   // true-已关注，false-未关注
	Avatar          string //用户头像
	BackgroundImage string //背景图像
	Signature       string //个人简介
	TotalFavorited  int64  //获赞数
	WorkCount       int64  //作品数
	FavoriteCount   int64  //点赞数
}

// TableName 修改表名映射
func (User) TableName() string {
	return "users"
}

// UserRequest 用户接口请求结构块
type UserRequest struct {
	Username string //注册用户名，最长32个字符
	Password string //密码，最长32个字符
}

// Response 用户响应状态码
type Response struct {
	StatusCode int32  `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty"` // 返回状态描述
}

// LoginResponse  登录接口响应结构快
type LoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"` // 用户id
	Token  string `json:"token"`             // 用户鉴权token
}

// UserResponse  用户接口响应结构快
type UserResponse struct {
	Response
	User *UserForUsed `json:"user"` //用户对象
}

// Token验证实例
var TokenInfo = map[string]User{
	"Token": {},
}
