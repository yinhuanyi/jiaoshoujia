/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-29
 * @Function:
 **/

package redisconnect

import "jiaoshoujia/models"

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
