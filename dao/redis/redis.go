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

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		fmt.Printf("Redis连接失败: %v\n", err)
		panic(err)
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
