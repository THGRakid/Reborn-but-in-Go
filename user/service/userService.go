package service

import (
	"Reborn-but-in-Go/user/dao"
	"Reborn-but-in-Go/user/model"
	"fmt"
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
func (s *UserService) CreateUser(username string, password string) (*model.LoginResponse, error) {

	// 调用 UserDao 的 UserLogin方法获取用户id及token
	user, token, _ := s.UserDao.CreateUser(username, password)

	if _, exist := model.TokenInfo[token]; exist {
		return &model.LoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户已存在"},
		}, nil
	} else {
		model.TokenInfo[token] = user

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
	userId, token, err := s.UserDao.UserLogin(username, password)
	if err != nil {
		// 处理错误，例如返回错误信息
		return nil, err
	}

	// 构建 UserResponse 对象，将查询到的消息记录填充进去
	loginResponse := &model.LoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	}
	return loginResponse, nil
}

// GetUserByID 根据用户ID和Token返回用户User列表
func (s *UserService) GetUserByID(userId int64) (*model.UserResponse, error) {
	user, _ := s.UserDao.GetUserByID(userId)

	userForUsed := &model.UserForUsed{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.Following,
		FollowerCount:   user.Follower,
		IsFollow:        false, //默认未关注
		Avatar:          user.Avatar,
		BackgroundImage: user.Background,
		Signature:       user.Introduce,
		TotalFavorited:  user.FavoritedCount,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	userResponse := &model.UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     userForUsed,
	}
	fmt.Println(userForUsed)
	return userResponse, nil
}
