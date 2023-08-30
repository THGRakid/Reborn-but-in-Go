package dao

import (
	"Reborn-but-in-Go/config"
	"Reborn-but-in-Go/user/model"
	"github.com/dgrijalva/jwt-go"
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
func (dao *UserDao) CreateUser(username string, password string) (model.User, string, error) {
	// 新建user类
	var user model.User
	/*
		// 检查用户名是否已存在
		existingUser := &model.User{}
		result := config.DB.Model(&model.User{}).Where("name = ?", username).First(existingUser)
		if result.Error == nil {
			return 0, "", errors.New("用户名已存在")
		}
	*/

	// 设置初始状态和创建时间
	user.Status = 0
	user.CreateAt = time.Now()

	//设置用户基本信息
	user.Name = username
	user.Password = password

	//将user内数据导入数据库
	config.DB.Create(&user)

	//创建Token并且绑定user
	token, _ := generateAuthToken(user.Id)

	// 创建成功后返回 user类型 和权限token
	return user, token, nil
}

// 生成权限token
func generateAuthToken(userID int64) (string, error) {
	// 检查 Redis 中是否已存在该用户的 token
	existingToken, err := config.RedisClient.Get(strconv.FormatInt(userID, 10)).Result()
	if err == nil && existingToken != "" {
		return existingToken, nil
	}

	// 生成 JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 设置 token 过期时间为一天

	// 签署 token
	jwtSecret := []byte("Reborn_but_in_Go") //暂时先设成硬编码
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// 存储生成的 token 到 Redis
	err = config.RedisClient.Set(strconv.FormatInt(userID, 10), tokenString, 0).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/*
方法二：
用户登录函数
参数：username string   用户名
参数：password string   密码
返回类型：userId，token, error
*/
func (dao *UserDao) UserLogin(username, password string) (int64, string, error) {
	// 根据用户名查询用户信息
	var user model.User

	result := config.DB.Table("user").Select("id, password").Where("name = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, "", errors.New("用户不存在")
		}
		return 0, "", result.Error
	}

	// 比较密码
	if user.Password == password {
		// 密码正确，返回用户ID
		token, _ := generateAuthToken(user.Id)
		return user.Id, token, nil
	}

	// 密码不正确，返回错误
	return 0, "", errors.New("密码不正确")
}

/*
方法三：
根据用户Token返回用户User指针
参数：userId string  用户Id
返回类型：*model.User （查询到的用户），error
*/

func (dao *UserDao) GetUserByID(userId string) (*model.User, error) {
	var user model.User

	result := config.DB.First(&user, userId)
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
