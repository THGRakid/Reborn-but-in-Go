package config

/*
import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"io"
	"os"
)

// 定义一个结构体来存储配置信息
type Config struct {
	Redis struct {
		Address  string `json:"IP地址"`
		Port     int    `json:"端口"`
		Password string `json:"密码"`
	} `json:"redis"`
}

func init() {
	// 打开配置文件
	configFile, err := os.Open("redis.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	// 在函数返回之前延迟执行 configFile.Close()，确保文件会被关闭
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {

		}
	}(configFile)

	// 读取配置文件的内容
	configBytes, _ := io.ReadAll(configFile)

	// 解析 JSON 格式的配置文件内容到 Config 结构体（反序列化）
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	// 配置 Redis 客户端
	ctx := context.Background()
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Address, config.Redis.Port), //拼接成完整的地址字符串
		Password: config.Redis.Password,                                         //提取密码
		DB:       0,                                                             //默认使用索引为 0 的数据库
	}
	client := redis.NewClient(options)

	// 测试连接
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

}


*/
