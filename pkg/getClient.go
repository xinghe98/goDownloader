package pkg

import (
	"net/http"
	"sync"
)

var (
	clientInstance *http.Client
	once           sync.Once
)

// 单例模式创建http请求
func GetClient() *http.Client {
	once.Do(func() {
		// 初始化逻辑
		clientInstance = &http.Client{
			// 初始化配置
		}
	})
	return clientInstance
}
