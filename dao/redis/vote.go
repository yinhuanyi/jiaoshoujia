/**
 * @Author: Robby
 * @File name: vote.go
 * @Create date: 2021-05-28
 * @Function:
 **/

package redisconnect

import (
	"errors"
	"time"
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

const (
	// 一天的秒数
	oneWeekInSeconds = 7 * 24 * 3600
)

func VoteForPost(userID, postID string, value float64) error {
	// 1：判断投票的限制
	// Redis Zscore 命令返回有序集中，成员的分数值，如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	// time.Now().Unix() 是获取当前时间戳的秒数, 返回的是int64类型
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	time.Now().Unix()
	// 2：更新帖子的分数
	// 3：记录用户为该帖子投票的数据
}
