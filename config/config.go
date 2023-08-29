package config

//sync 是 Go 语言标准库中提供的用于并发操作的包。
//它包含了一些用于同步的原语，如互斥锁（Mutex）和读写锁（RWMutex）
import (
	"sync"
)

var (
	serverIP          = ""
	serverAddress     = "http://localhost:8080" // 默认服务器地址
	serverAddressLock sync.RWMutex              // 用于保护服务器地址并发访问
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
