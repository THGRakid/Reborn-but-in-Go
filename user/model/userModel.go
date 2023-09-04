package model

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
func (User) TableName() string {
	return "users"
}

// Response 用户响应状态码
type Response struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

// LoginResponse  登录接口响应结构快
type LoginResponse struct {
	Response
	UserId int64  `json:"user_id"` // 用户id
	Token  string `json:"token"`   // 用户鉴权token
}

// UserResponse  用户接口响应结构快
type UserResponse struct {
	Response
	User *User `json:"user"` //用户对象
}
