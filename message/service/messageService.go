package service

import (
	"Reborn-but-in-Go/message/dao"
	"Reborn-but-in-Go/message/model"
	"time"
)

// MessageService 服务层
type MessageService struct {
	MessageDao *dao.MessageDao
}

// NewMessageService 创建一个新的 MessageService 实例
func NewMessageService(messageDao *dao.MessageDao) *MessageService {
	return &MessageService{
		MessageDao: messageDao,
	}
}

// 假定一个用户id，这个应该从user方法获取
var UserId int64 = 123

// QueryMessage 根据用户ID和上次消息时间获取聊天消息记录
func (s *MessageService) QueryMessage(toUserId int64, preMsgTime string) (*model.ChatResponse, error) {

	// 调用 MessageDao 的 QueryMessageList 方法获取聊天消息记录

	messageList := s.MessageDao.QueryMessage(&preMsgTime, UserId, toUserId)

	// 构建 ChatResponse 对象，将查询到的消息记录填充进去
	chatResponse := &model.ChatResponse{
		StatusCode:  0,
		StatusMsg:   "Success",
		MessageList: messageList,
	}
	return chatResponse, nil
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(toUserID int64, content string) error {
	message := &model.Message{
		UserId:   UserId,
		ToUserId: toUserID,
		Content:  content,
		CreateAt: time.Now(),
		Status:   1,
	}

	// 调用 DAO 的 CreateMessage 方法来保存消息到数据库
	err := s.MessageDao.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}
