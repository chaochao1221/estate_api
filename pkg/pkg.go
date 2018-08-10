package pkg

import (
	"estate/pkg/redis"
)

// 初始化服务器组件和sdk
func Init() {
	redis.NewRedisCache()
}
