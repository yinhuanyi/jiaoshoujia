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

	"github.com/go-redis/redis"
)

func createScript() *redis.Script {
	script := redis.NewScript(`
		
			-- 声明局部变量
			local goodsSurplus
			local flag
			
			-- existUserIds：hadBuyUids
			-- goodsSurplusKey：goodsSurplus
			local existUserIds    = tostring(KEYS[1])
			local goodsSurplusKey = tostring(KEYS[2])
			
			-- memberUid：5824742984
			local memberUid       = tonumber(ARGV[1])
			
			-- 执行Redis命令：查看existUserIds集合中是否存在memberUid
			local hasBuy = redis.call("sIsMember", existUserIds, memberUid)
			
			-- 判断用户是否买了商品，如果为0，说明没有买，如果不为0，说明买了，那么直接结束脚本
			if hasBuy ~= 0 then
			  return 0
			end
			
			-- 如果用户没有买，获取goodsSurplusKey商品的库存，如果库存不存在，那么直接结束脚本
			goodsSurplus =  redis.call("GET", goodsSurplusKey)
			if goodsSurplus == false then
			  return 0
			end
			
			-- 如果库存存在，且库存为0，那么直接结束脚本
			goodsSurplus = tonumber(goodsSurplus)
			if goodsSurplus <= 0 then
			  return 0
			end
			
			-- 扣减库存：先将用户的id加入到已买集合，再让库存减一
			flag = redis.call("SADD", existUserIds, memberUid)
			flag = redis.call("DECR", goodsSurplusKey)
			
			-- 最后返回1，表示扣减库存成功
			return 1
	`)
	return script
}

// evalScript 执行lua脚本
func evalScript(client *redis.Client, userId string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取脚本对象
	script := createScript()

	// 创建sha值与脚本的对应关系
	sha, err := script.Load(client).Result()
	if err != nil {
		log.Fatalln(err)
	}

	// 执行sha值对应的脚本，key是"hadBuyUids", "goodsSurplus"，args是userId，例如5824742984，keys和args会传递到Redis脚本中
	ret := client.EvalSha(sha, []string{"hadBuyUids", "goodsSurplus"}, userId)

	// 获取Redis的返回结果
	if result, err := ret.Result(); err != nil {
		log.Fatalf("Execute Redis fail: %v\n", err.Error())
	} else {
		if result.(int64) == 1 {
			fmt.Printf("userId：%s，第一次抢购商品，扣减库存成功\n", userId)
		} else {
			fmt.Printf("userId：%s，已经抢购过商品，扣减库存失败\n", userId)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 这里假设4个用户发送了6个请求去扣减库存
	for _, v := range []string{"5824742984", "5824742984", "5824742983", "5824742983", "5824742982", "5824742980"} {
		wg.Add(1)
		go evalScript(client, v, &wg)
	}
	wg.Wait()

}
