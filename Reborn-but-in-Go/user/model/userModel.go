package model

import (
	"time"
)

// Message 持久层结构块
type User struct {
	Id        	int64     `gorm:"primaryKey"` //用户id
	Name    	string    `gorm:"unique"` //用户名称（不重复）
	Follower 	int64      //粉丝人数
	Following   int64      //关注人数
	Password  	string 	   //密码
	Avatar      string     //用户头像
	Background  string      //背景图像
	Introduce   string      //个人简介
	Favorited_count int64  //获赞数
	Work_count int64 		//作品数
	Favorite_count int64    //点赞数
	CreateAt     	time.Time	//用户创建时间
	Status      int8        //用户状态（在线为1，不在线为0）
}

// TableName 修改表名映射
func (User) TableName() string {
	return "user"
}
