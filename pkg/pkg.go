package pkg

import (
	"k12_marking_api/pkg/redis"
)

// 初始化服务器组件和sdk
func Init() {
	redis.NewRedisCache()
}
