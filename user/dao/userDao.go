package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/user/model"
	"crypto/rand"
	"encoding/base64"

	"gorm.io/gorm"
	"strconv"
	//"gorm.io/gorm/logger"
	//"fmt"
	"errors"
	"sync"
	"time"
)

/*
	构造方法
*/
// UserDao 用于操作消息数据的数据库访问对象（DAO）
type UserDao struct {
}

// 用于保存单例实例
var userDao *UserDao

// 单例模式，只生成一个userDao实例，提高性能
// sync.Once 是一个同步原语（synchronization primitive）
// 用于确保在并发环境下只执行一次特定的操作。它通常用于延迟初始化或只需要在程序的生命周期内执行一次的操作
var userOnce sync.Once

// NewUserDaoInstance 用于获取 UserDao 单例实例的函数
// 传递一个匿名函数（闭包），其作用是创建一个新的 UserDao 实例并将其赋值给 userDao变量。
// 这个函数只会在第一次调用 Do 方法时执行
func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

/*
方法一：
创建用户
参数：username string用户名, password string密码
返回值：id int 用户id，token string 令牌，error错误
*/
func (dao *UserDao) CreateUser(username string, password string) (int64, string, error) {
	var user model.User
	// 检查用户名是否已存在
	existingUser := &model.User{}
	result := config.DB.Model(&model.User{}).Where("name = ?", username).First(existingUser)
	if result.Error == nil {
		return 0, "", errors.New("用户名已存在")
	}

	// 设置初始状态和创建时间
	user.Status = 0
	user.CreateAt = time.Now()

	//设置用户基本信息
	user.Name = username
	user.Password = password

	//将user内数据导入数据库
	result = config.DB.Create(&user)
	if result.Error != nil {
		return 0, "", result.Error
	}

	// 创建成功后返回用户 id 和权限token
	temp_token, _ := generateAuthToken(user.Id)
	return user.Id, temp_token, nil
}

// 生成权限token
func generateAuthToken(userID int64) (string, error) {
	// 生成一个随机的字节数组作为令牌
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// 将字节数组进行Base64编码，生成字符串作为令牌
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	// 将userID与令牌进行组合，以确保唯一性
	userIDString := strconv.FormatInt(userID, 10)
	token = userIDString + "-" + token

	return token, nil
}

/*
方法二：
用户登录函数
参数：username string   用户名
参数：password string   密码
返回类型：userId，error
*/
func (dao *UserDao) UserLogin(username, password string) (int64, error) {
	// 根据用户名查询用户信息
	var user model.User

	result := config.DB.Table("user").Select("id, password").Where("name = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("用户不存在")
		}
		return 0, result.Error
	}

	// 比较密码
	if user.Password == password {
		// 密码正确，返回用户ID
		return user.Id, nil
	}

	// 密码不正确，返回错误
	return 0, errors.New("密码不正确")
}

/*
方法三：
根据用户ID和Token返回用户User指针
参数：userID int64  用户名
返回类型：*model.User （查询到的用户），error
*/

func (dao *UserDao) GetUserByID(userID int64) (*model.User, error) {
	var user model.User
	result := config.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

/*
已废弃
方法四：获取个人关注列表
根据用户ID查看用户信息
参数：userID int64  用户名
返回类型：*model.User （查询到的用户），error

func GetFollowing(userID int64) ([]int64, error) {
	var user model.User
	result := config.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return
}*/
