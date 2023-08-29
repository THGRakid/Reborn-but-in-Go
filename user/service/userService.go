package service

import (
	"Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/user/model"
)

// UserService 服务层
type UserService struct {
	UserDao *dao.UserDao
	//RedisClient *redis.Client
}

// NewUserService 创建一个新的 UserService 实例
func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{
		UserDao: userDao,
	}
}

type IdResponse struct {
	StatusCode int32       // 状态码，0-成功，其他值-失败
	StatusMsg  string      // 返回状态描述
	User       *model.User //用户信息
}

// CreateUser 根据用户名和登录密码注册用户id及token
func (s *UserService) CreateUser(username string, password string) (*model.UserResponse, error) {

	// 调用 UserDao 的 UserLogin方法获取用户id及token

	userId, token, _ := s.UserDao.CreateUser(username, password)

	// 构建 UserResponse 对象，将查询到的消息记录填充进去
	userResponse := &model.UserResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		UserId:     userId,
		Token:      token,
	}
	return userResponse, nil
}

// UserLogin 根据用户名和登录密码获取用户id及token
func (s *UserService) UserLogin(username string, password string) (*model.UserResponse, error) {

	// 调用 UserDao 的 UserLogin 方法获取用户信息
	userId, err := s.UserDao.UserLogin(username, password)
	if err != nil {
		// 处理错误，例如返回错误信息
		return nil, err
	}

	token := "1" // 需修改为生成的实际 token
	// 构建 UserResponse 对象，将查询到的消息记录填充进去
	userResponse := &model.UserResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		UserId:     userId,
		Token:      token,
	}
	return userResponse, nil
}

// GetUserByID 根据用户ID和Token返回用户User列表
func (s *UserService) GetUserByID(userId int64) (*IdResponse, error) {
	user, _ := s.UserDao.GetUserByID(userId)

	idResponse := &IdResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		User:       user,
	}
	return idResponse, nil
}
