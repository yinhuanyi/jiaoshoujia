/**
 * @Author: Robby
 * @File name: vote.go
 * @Create date: 2021-05-28
 * @Function:
 **/

/*
投一票就加432分，86400/200, 那么天时间是86400，200张赞成票的帖子，才能在网站上置顶

direction = 1 时，有两种情况：
	1：之前没有投过票，现在投赞成票：|0-1| = 1   分数：432*1
	2：之前投反对票，现在改投赞成票：|-1-1| = 2  分数：432*2
	3：之前是赞成票，现在改赞成票：  |1-1| = 0
direction = 0 时，有两种情况：
	1：之前投过赞成票，现在取消投票：|1-0|=1      分数：432*-1
	2：之前投反对票，现在取消投票：  |-1-0|=1     分数：432*1

direction = -1 时，有两种情况：
	1：之前没有投过票，现在投反对票：|0--1|=1      分数：432*-1
	2：之前投赞成票，现在改反对票：  |1--1|=2      分数：432*-2

投票的限制：
	每个帖子一周以后就不能再次投票了
	到期之后，将redis中保存的赞成票数和反对票数存储到MySQL数据库
	到期之后，删除redis中keyPostVotedZSetPF这个key的分数
*/

package redisconnect

import (
	"errors"
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

const (
	// 一天的秒数
	oneWeekInSeconds = 7 * 24 * 3600
	// 每一票是 432分
	scorePerVote = 432
)

// VoteForPost  value表示用户投票是1、0、-1
func VoteForPost(userID, postID string, value float64) (err error) {
	// (一)：判断帖子允许投票的时间
	// 这里的postTime是获取帖子的【发帖时间】，redis获取, 这里是把这个Redis的分数，作为发帖时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	// time.Now().Unix() 是获取当前时间戳的秒数, 返回的是int64类型
	// 如果发帖的时间距离现在的时间大于一天，返回err
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		err = ErrVoteTimeExpire
	}

	// (二)：更新帖子的分数
	// 将更新分数和记录用户为该帖子投票的数据作为一个原子的事务操作
	pipeline := client.TxPipeline()

	// 查看之前的投票记录, keyPostVotedZSetPF+postID这个key是返回userId的分数，作为某个用户是否为某个帖子投过票， 如果是1表示投过赞成票，如果是-1表示投过反对票，如果是0表示没有投过票
	ov := client.ZScore(getRedisKey(keyPostVotedZSetPF+postID), userID).Val()

	// 这op是存储方向
	var op float64
	// 这里是计算分数的正负号
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	// 计算绝对值
	diff := math.Abs(ov - value)

	// 获取当前这个帖子的分数
	OldpostScore := client.ZScore(getRedisKey(keyPostScoreZSet), postID).Val()
	// 计算本次用户给帖子的分数
	CurpostScore := op * diff * scorePerVote

	// 更新帖子的分数(原子操作)
	pipeline.ZIncrBy(getRedisKey(keyPostScoreZSet), OldpostScore+CurpostScore, postID)

	// (三)：记录用户为该帖子投票的数据
	if value == 0 {
		// 如果用户是取消投票，那么直接删除用户为此帖子的投票记录
		pipeline.ZRem(getRedisKey(keyPostVotedZSetPF+postID), userID)
	} else {
		// 如果不是取消投票，是赞成或反对票，那么直接添加一个元素，如果此元素之前存在，那么就会更新分数
		pipeline.ZAdd(getRedisKey(keyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec() // 执行事务
	return
}

// CreatePostTime 在redis中，给新帖子创建时间
func CreatePostTimeAndScore(postId int64) (err error) {

	zap.L().Debug("client.ZAdd(getRedisKey(KeyPostTimeZSet)", zap.Int64("postId", postId))

	// 让帖子的时间初始化和分数初始化作为一个原子的事务操作执行
	pipeline := client.TxPipeline() // 申请一个事务

	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	pipeline.ZAdd(getRedisKey(keyPostScoreZSet), redis.Z{
		Score:  float64(0),
		Member: postId,
	})

	_, err = pipeline.Exec() // 执行Redis原子事务操作

	return
}
