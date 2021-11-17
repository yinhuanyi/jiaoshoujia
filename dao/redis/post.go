/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-29
 * @Function:
 **/

package redisconnect

import (
	"jiaoshoujia/models"

	"github.com/go-redis/redis"
)

// GetPostVoteData ：根据post的ids，查询每一篇帖子投赞成票的分数, 注意这里没有获取未投票或者
func GetPostVoteData(ids []string) (data []int64, err error) {

	// Todo: 1：由于每次遍历，都会连接一次redis，导致redis连接开销较大，这里可以优化一下，改为将所有的key全部发给redis，再一次性获取结果
	//for _, id := range ids {
	//	key := getRedisKey(keyPostVotedZSetPF + id)
	//	// 返回 每个key对应的分数中，分数值为1的成员个数，也就是投赞成票的赞成票的成员个数，也就是投赞成票的用户个数
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}

	// Todo: 2：这是使用redis的pipeline来解决, 一次发送多条命令，获取一条结果，减少RTT
	pipeline := client.Pipeline() // 声明redis的pipeline
	for _, id := range ids {
		key := getRedisKey(keyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1") // 这里只是添加，没有执行
	}

	cmders, err := pipeline.Exec() // 这里统一执行
	if err != nil {
		return
	}

	for _, cmder := range cmders {
		// 类型断言，类型转换一下
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

// GetPostIdsInOrder 排序
func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1：根据用户传递的order参数，判断redis中的key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(keyPostScoreZSet)
	}

	// 2：确定查询索引的起始点, 因为传递了分页的参数
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	// 3：按照元素的分数，从大到小查询指定数量的元素，返回一个string类型的列表
	return client.ZRevRange(key, start, end).Result()

}
