package dao

import (
	"fmt"
	"sync"
	"time"
	"字节青训营/Reborn-but-in-Go/config"
	"字节青训营/Reborn-but-in-Go/message/model"
)

/*
	构造方法
*/
// MessageDao 用于操作消息数据的数据库访问对象（DAO）
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

/*
方法一：
创建一条消息
参数：message对象，里面包含：send_id，receive_id，content
*/
func (*MessageDao) CreateMessage(message *model.Message) error {
	// 设置初始状态和创建时间
	message.Status = 1
	message.CreateAt = time.Now()

	//将message内数据导入数据库
	result := config.DB.Create(&message)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

/*
方法二：
查看在一段时间内双方的聊天记录
参数：date string   一段时间
参数：SendId int64   此信息发送者
参数：ReceiveId int64   此信息接受者
返回值：message对象列表
*/
func (*MessageDao) QueryMessageList(date *string, SendId int64, ReceiveId int64) []*model.Message {
	fmt.Println(*date)
	var MessageList []*model.Message
	//查询双方信息，并指定时间范围，按从小到大顺序
	config.DB.Where("( (send_id = ? and receive_id = ?) or (send_id = ? receive_id = ?) ) and create_at > ?", SendId, ReceiveId, ReceiveId, SendId, date).
		Order("create_at asc").Find(&MessageList)

	fmt.Println(MessageList)
	return MessageList
}

/*

*

func (d *MessageDao) QueryMessage(toUserId int64, FromUserId int64) *to_relation.QueryBody {
	message := Message{}
	var msgType int64
	err := DB.Model(&Message{}).Where("(to_user_id=? and from_user_id=?) or (to_user_id=? and from_user_id=?)", toUserId, FromUserId, FromUserId, toUserId).Order("create_at desc").First(&message).Error
	if err != nil { //没查到，first会报个错
		return &to_relation.QueryBody{
			FromUserId: FromUserId,
			ToUserId:   toUserId,
			Message:    &to_relation.Message{},
			MsgType:    0,
		}
	}
	if FromUserId == message.FromUserId {
		msgType = 1
	} else {
		msgType = 0
	}

	return &to_relation.QueryBody{
		FromUserId: FromUserId,
		ToUserId:   toUserId,
		Message: &to_relation.Message{
			Id:      message.Id,
			Content: message.Content,
		},
		MsgType: msgType,
	}
}
*/
