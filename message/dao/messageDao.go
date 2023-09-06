package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/message/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
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

// QueryMessage
/*
  方法一：
  查看在一段时间内双方的聊天记录
  参数：date string   一段时间
  参数：UserId int64   此信息发送者
  参数：ToUserId int64   此信息接受者
  返回值：message对象列表
*/
func (*MessageDao) QueryMessage(date string, UserId int64, ToUserId int64) []*model.Message {
	fmt.Println(date)
	fmt.Println(UserId, "这是UserId")
	fmt.Println(ToUserId, "这是ToUserId")
	time.Sleep(5 * time.Second)
	var MessageList []*model.Message
	//查询双方信息，并指定时间范围，按从小到大顺序
	config.DB.Where("((user_id = ? and to_user_id = ?) or (user_id = ? and to_user_id = ?))", UserId, ToUserId, ToUserId, UserId).
		Order("create_at asc").Find(&MessageList)

	fmt.Println(MessageList, "这是MessageList")
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

// QueryMessageList
/*
  方法三：
  获取列表最新消息
  参数：id int64   与信息有关的人
  返回值：message对象列表
*/
func (*MessageDao) QueryMessageList(id int64) []model.FriendUser {
	var MessageList []model.FriendUser

	// 查询当前请求用户发送的消息
	sentMessages := []struct {
		Content  string
		UserId   int64
		ToUserId int64
		CreateAt time.Time
		UserInfo model.User
		MsgType  int64
	}{}
	config.DB.Model(&model.Message{}).
		Select("content, to_user_id,create_at").
		Where("user_id = ? ", id).
		Order("create_at ASC").
		Find(&sentMessages)

	// 发送消息 MsgType 为1
	for i := range sentMessages {
		sentMessages[i].MsgType = 1
	}

	// 遍历查询结果并为每个 to_user_id 调用 getUserById 方法
	for i := range sentMessages {
		user, _ := GetUserById(sentMessages[i].ToUserId)
		sentMessages[i].UserInfo = user
	}

	// 查询当前请求用户接收的消息
	receivedMessages := []struct {
		Content  string
		UserId   int64
		ToUserId int64
		CreateAt time.Time
		UserInfo model.User
		MsgType  int64
	}{}
	config.DB.Model(&model.Message{}).
		Select("content, user_id,create_at").
		Where("to_user_id = ? ", id).
		Order("create_at ASC").
		Find(&receivedMessages)

	// 接收消息 MsgType 为1
	for i := range receivedMessages {
		receivedMessages[i].MsgType = 0
	}

	// 遍历查询结果并为每个 user_id 调用 getUserById 方法
	for i := range receivedMessages {
		user, _ := GetUserById(receivedMessages[i].UserId)
		receivedMessages[i].UserInfo = user
	}

	// 删除l两个人多次消息，保留 createAt 较晚的那条消息
	cleanedSentMessages := []struct {
		Content  string
		UserId   int64
		ToUserId int64
		CreateAt time.Time
		UserInfo model.User
		MsgType  int64
	}{}

	for _, sentMessage := range sentMessages {
		var shouldKeep = true

		for _, receivedMessage := range receivedMessages {
			if sentMessage.ToUserId == receivedMessage.UserId {
				// 如果 userId 和 toUserId 相等，比较 createAt
				if sentMessage.CreateAt.After(receivedMessage.CreateAt) {
					// 如果 sentMessage 的 createAt 较晚，保留它
					shouldKeep = true
				} else {
					// 否则，删除 sentMessage
					shouldKeep = false
				}
				break // 已经找到匹配的 receivedMessage，退出内部循环
			}
		}

		if shouldKeep {
			cleanedSentMessages = append(cleanedSentMessages, sentMessage)
		}
	}

	// 将 receivedMessages 中未匹配的消息也加入 cleanedSentMessages
	for _, receivedMessage := range receivedMessages {
		var found = false

		for _, sentMessage := range cleanedSentMessages {
			if sentMessage.ToUserId == receivedMessage.UserId {
				found = true
				break
			}
		}

		if !found {
			cleanedSentMessages = append(cleanedSentMessages, receivedMessage)
		}
	}

	// 提取内容组成FriendUser并添加到MessageList中
	for _, message := range cleanedSentMessages {
		MessageList = append(MessageList, model.FriendUser{
			Message:         message.Content,
			MsgType:         message.MsgType,
			Id:              message.UserInfo.Id,
			Name:            message.UserInfo.Name,
			FollowCount:     message.UserInfo.FollowCount,
			FollowerCount:   message.UserInfo.FollowerCount,
			IsFollow:        message.UserInfo.IsFollow,
			Avatar:          message.UserInfo.Avatar,
			BackgroundImage: message.UserInfo.BackgroundImage,
			Signature:       message.UserInfo.Signature,
			TotalFavorited:  message.UserInfo.TotalFavorited,
			WorkCount:       message.UserInfo.WorkCount,
			FavoriteCount:   message.UserInfo.FavoriteCount,
		})
	}

	fmt.Println(MessageList)
	return MessageList
}

// GetUserById
/*
方法三：
根据用户Token返回用户User指针
参数：userId string  用户Id
返回类型：*model.User （查询到的用户），error
*/

func GetUserById(id int64) (model.User, error) {
	var user model.User

	result := config.DB.Table("users").Where("id = ?", id).Limit(1).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, errors.New("用户不存在")
		}
		return user, result.Error
	}
	return user, nil
}
