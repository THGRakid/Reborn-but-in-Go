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

// CreateUser 根据用户名和登录密码注册用户id及token
func (s *UserService) CreateUser(username string, password string) (*model.LoginResponse, error) {

	// 调用 UserDao 的 UserLogin方法获取用户id及token
	user, token, _ := s.UserDao.CreateUser(username, password)

	checkInfo, err := s.UserDao.CheckUser(username)
	if checkInfo == 1 {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 2, StatusMsg: "用户已存在，请更改用户名"},
		}, err
	} else {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		}, nil
	}

}

// UserLogin 根据用户名和登录密码获取用户id及token
func (s *UserService) UserLogin(username string, password string) (*model.LoginResponse, error) {

	// 调用 UserDao 的 UserLogin 方法获取用户信息
	userId, token, checkInfo, err := s.UserDao.UserLogin(username, password)
	if err != nil {
		// 处理错误，例如返回错误信息
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 3, StatusMsg: "登录错误"},
		}, err
	}

	if checkInfo == 1 {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 2, StatusMsg: "用户不存在"},
		}, err
	} else if checkInfo == 2 {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 3, StatusMsg: "登录错误"},
		}, err
	} else if checkInfo == 3 {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 2, StatusMsg: "密码错误"},
		}, err
	} else if checkInfo == 0 {
		// 构建 UserResponse 对象，将查询到的消息记录填充进去
		loginResponse := &model.LoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   userId,
			Token:    token,
		}
		return loginResponse, nil
	} else {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 3, StatusMsg: "登录错误"},
		}, err
	}

}

// GetUserByID 根据用户ID和Token返回用户User列表
func (s *UserService) GetUserByID(userId int64) (*model.UserResponse, error) {
	user, _ := s.UserDao.GetUserByID(userId)

	userResponse := &model.UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     user,
	}

	return userResponse, nil
}
