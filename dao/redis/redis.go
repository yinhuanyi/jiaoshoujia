/**
 * @Author: Robby
 * @File name: redis.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package redisconnect

import (
	"fmt"
	"jiaoshoujia/settings"
	"log"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var client *redis.Client

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v\n", err)
	}

	return nil
}

func Close() {
	_ = client.Close()
}
