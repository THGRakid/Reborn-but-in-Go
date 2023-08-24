package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/message/model"
	"fmt"
	"sync"
)

/*
	构造方法
*/
// MessageDao 数据层
type MessageDao struct {
}

// 用于保存单例实例
var messageDao *MessageDao

// 单例模式，只生成一个messageDao实例，提高性能
// sync.Once 是一个同步原语（synchronization primitive）
// 用于确保在并发环境下只执行一次特定的操作。它通常用于延迟初始化或只需要在程序的生命周期内执行一次的操作
var messageOnce sync.Once

// NewMessageDaoInstance 用于获取 MessageDao 单例实例的函数
// 传递一个匿名函数（闭包），其作用是创建一个新的 MessageDao 实例并将其赋值给 messageDao 变量。
// 这个函数只会在第一次调用 Do 方法时执行
func NewMessageDaoInstance() *MessageDao {
	messageOnce.Do(
		func() {
			messageDao = &MessageDao{}
		})
	return messageDao
}

// QueryMessageList
/*
  方法一：
  查看在一段时间内双方的聊天记录
  参数：date string   一段时间
  参数：SendId int64   此信息发送者
  参数：ReceiveId int64   此信息接受者
  返回值：message对象列表
*/
func (*MessageDao) QueryMessage(date *string, UserId int64, ToUserId int64) []*model.Message {
	fmt.Println(*date)
	var MessageList []*model.Message
	//查询双方信息，并指定时间范围，按从小到大顺序
	config.DB.Where("( (send_id = ? and receive_id = ?) or (send_id = ? receive_id = ?) ) and create_at > ?", UserId, ToUserId, ToUserId, UserId, date).
		Order("create_at asc").Find(&MessageList)

	fmt.Println(MessageList)
	return MessageList
}

// CreateMessage
/*
  方法二：
  创建一条消息
  参数：message对象
*/
func (*MessageDao) SendMessage(message *model.Message) error {
	//将message内数据导入数据库
	result := config.DB.Create(&message)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

/*
	方法三：
	展示页面聊天记录（最后一条）
*/
