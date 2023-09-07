package config

// 使用云服务器
//用户名：administrator
//密码：Go123456@++
import (
	"sync"
)

var (
	serverAddress     = "http://l47.113.230.49:8090" // 默认服务器地址
	serverAddressLock sync.RWMutex                   // 用于保护服务器地址并发访问
)

// SetServerAddress 设置服务器地址
func SetServerAddress(address string) {
	serverAddressLock.Lock()
	defer serverAddressLock.Unlock()

	serverAddress = address
}

// GetServerAddress 获取服务器地址
func GetServerAddress() string {
	serverAddressLock.RLock()
	defer serverAddressLock.RUnlock()

	return serverAddress
}
