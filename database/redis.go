package database

import (
	"FullTimeTeacher/config"
	"FullTimeTeacher/log/logx"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 全局redis客户端
var RDB *redis.Client

const (
	Redis_Token_Key             = "token:%s"
	Redis_Captcha_Key           = "captcha:%s"
	Redis_Verification_Code_Key = "verification_code:%s"
)

// 初始化redis客户端
func InitRedis(config config.RedisConfig) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		logx.GetLogger("logx").Errorf("Redis连接失败: %v", err)
	} else {
		logx.GetLogger("logx").Info("Redis连接成功")
	}
}
