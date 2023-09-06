package service

import (
	"Reborn-but-in-Go/message/dao"
	"Reborn-but-in-Go/message/model"
	"fmt"
	"time"
)

// MessageService 服务层
type MessageService struct {
	MessageDao *dao.MessageDao
	//RedisClient *redis.Client
}

// NewMessageService 创建一个新的 MessageService 实例
func NewMessageService(messageDao *dao.MessageDao) *MessageService {
	return &MessageService{
		MessageDao: messageDao,
	}
}

// QueryList 根据用户ID和上次消息时间获取聊天消息记录
func (s *MessageService) QueryList(UserId int64) (*model.ListResponse, error) {

	// 调用 MessageDao 的 QueryMessageList 方法获取聊天消息记录

	messageList := s.MessageDao.QueryMessageList(UserId)

	// 构建 ListResponse 对象，将查询到的消息记录填充进去
	listResponse := &model.ListResponse{
		UserList: messageList,
	}
	return listResponse, nil
}

// QueryMessage 根据用户ID和上次消息时间获取聊天消息记录
func (s *MessageService) QueryMessage(UserId int64, toUserId int64, preMsgTime string) (*model.ChatResponse, error) {

	// 调用 MessageDao 的 QueryMessageList 方法获取聊天消息记录

	messageList := s.MessageDao.QueryMessage(preMsgTime, UserId, toUserId)

	// 构建 ChatResponse 对象，将查询到的消息记录填充进去
	chatResponse := &model.ChatResponse{
		MessageList: messageList,
	}
	return chatResponse, nil
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(userId int64, toUserId int64, content string) error {
	message := &model.Message{
		UserId:   userId,
		ToUserId: toUserId,
		Content:  content,
		CreateAt: time.Now(),
	}

	fmt.Println(message, "这是我")
	// 调用 DAO 的 CreateMessage 方法来保存消息到数据库
	err := s.MessageDao.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}
