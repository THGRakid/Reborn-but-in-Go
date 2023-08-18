package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/user/model"

	//"gorm.io/driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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
参数：user对象，里面包含：send_id，receive_id，content
*/
func (*UserDao) CreateUser(user *model.User) error {
	// 设置初始状态和创建时间
	user.Status = 0
	user.CreateAt = time.Now()

	//将user内数据导入数据库
	result := config.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

/*
方法二：
用户登录函数
参数：username string   用户名
参数：password string   密码
返回类型：bool （是否创建成功），error
*/

func UserLogin(username, password string) (bool, error) {
	// 根据用户名查询用户信息
	var user model.User
	result := config.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("用户不存在")
		}
		return false, result.Error
	}

	// 验证密码是否匹配
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, errors.New("密码不正确")
	}

	// 登录成功
	return true, nil
}

/*
方法三：
根据用户ID查看用户信息
参数：userID int64  用户名
返回类型：*model.User （查询到的用户），error
*/

func GetUserByID(userID int64) (*model.User, error) {
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
