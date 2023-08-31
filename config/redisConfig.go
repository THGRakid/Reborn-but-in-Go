package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"os"
	"path/filepath"
)

// 存储从 JSON 配置文件中读取的 Redis 配置信息
type RedisConfig struct {
	Redis struct {
		Address  string `json:"IP地址"`
		Port     int    `json:"端口"`
		Password string `json:"密码"`
	} `json:"redis"`
}

// 可以调用的客户端
var RedisClient *redis.Client

// 配置并创建 Redis 客户端连接
func setupRedisClient(config RedisConfig) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port),
		Password: "", //设置无密码
		DB:       0,
	}
	client := redis.NewClient(options)

	// 连接测试
	pong, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Redis 打开错误: %v", err)
	}
	fmt.Println("正在连接 Redis :", pong)

	return client, nil
}

func InitRedis() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前工作目录错误:", err)
		return
	}
	// 构建 redis.json 的完整路径
	configFilePath := filepath.Join(currentDir, "\\config\\redis.json")

	// 打开 redis.json 配置文件
	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("配置文件打开错误:", err)
		return
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println("配置文件关闭错误:", err)
		}
	}(configFile)

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Println("配置文件读取错误:", err)
		return
	}

	var config RedisConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}

	RedisClient, err = setupRedisClient(config)
	if err != nil {
		fmt.Println("初始化 Redis 客户端错误:", err)
		return
	}

	fmt.Println("Redis 客户端初始化成功")
}
