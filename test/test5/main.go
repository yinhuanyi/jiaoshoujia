/**
 * @Author: Robby
 * @File name: main.go
 * @Create date: 2021-10-21
 * @Function:
 **/

package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/cch123/goroutineid"

	"github.com/go-redis/redis"
)

var wg sync.WaitGroup

var client *redis.Client

type ReentrantLock struct{}

func NewReentrantLock() sync.Locker {
	return &ReentrantLock{}
}

// lockScript 获取锁脚本
func lockScript() *redis.Script {
	script := redis.NewScript(`
		-- key：lock
		local key = KEYS[1]
		local goroutineId = ARGV[1]
		-- timeout: 20秒
		local timeout = ARGV[2]
		
		-- 判断lock这个key不存在
		if (redis.call('exists', key) == 0) then
			-- 获取锁，记录重入次数，第一次为1
			redis.call('hset', key, goroutineId, '1')
			-- 设置超时时间
			redis.call('expire', key, timeout)
			return 1
		end
		
		-- 判断是不是自己的锁，如果是自己的锁
		if (redis.call('hexists', key, goroutineId) == 1) then
			-- 直接加1
			redis.call('hincrby', key, goroutineId, '1')
			-- 设置超时时间
			redis.call('expire', key, timeout)
			return 1
		end
		
		-- 如果是别人的锁，说明获取锁失败，返回0
		return 0

	`)
	return script
}

// unlockScript 释放锁脚本
func unlockScript() *redis.Script {
	script := redis.NewScript(`
		-- key：lock
		local key = KEYS[1]
		local goroutineId = ARGV[1]
		
		-- 如果不存在，直接返回0
		if (redis.call('HEXISTS', key, goroutineId) == 0) then
			return 1
		end
		
		-- 判断重入次数是否为0，如果为0，说明锁全部释放，需要删除锁
		local count = redis.call('HINCRBY', key, goroutineId, -1)
		
		if (count == 0) then
			redis.call('DEL', key)
			return 1
		end

		return 1

	`)
	return script
}

// Lock 加锁
func (r *ReentrantLock) Lock() {

	lock := lockScript()
	sha, err := lock.Load(client).Result()
	if err != nil {
		panic(err)
	}
	// 执行脚本
	goroutineId := goroutineid.GetGoID()

	ret := client.EvalSha(sha, []string{"lock"}, goroutineId, 20)

	if result, err := ret.Result(); err != nil {
		log.Panicf("获取锁Execute lua fail: %v\n", err.Error())
	} else {

		if result.(int64) == 1 {
			fmt.Printf("goroutineId：%d，获取锁成功\n", goroutineId)
		} else {
			fmt.Printf("goroutineId：%d，获取锁失败\n", goroutineId)
		}
	}
}

// Unlock 解锁
func (r *ReentrantLock) Unlock() {

	unlock := unlockScript()
	sha, err := unlock.Load(client).Result()
	if err != nil {
		panic(err)
	}
	// 执行脚本
	goroutineId := goroutineid.GetGoID()
	ret := client.EvalSha(sha, []string{"lock"}, goroutineId)

	if _, err = ret.Result(); err != nil {
		log.Panicf("释放锁Execute lua fail: %v\n", err.Error())
	}

	fmt.Printf("goroutineId：%d，释放锁成功\n", goroutineId)
}

func main() {

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	lock := NewReentrantLock()

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			lock.Lock()
			lock.Lock()

			lock.Unlock()
			lock.Unlock()
			lock.Unlock()

		}()
	}

	wg.Wait()

	fmt.Println("执行完毕")

}
