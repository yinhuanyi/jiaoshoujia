/**
 * @Author: Robby
 * @File name: vote.go
 * @Create date: 2021-05-27
 * @Function:


 **/

package logic

import (
	redisconnect "jiaoshoujia/dao/redis"
	"jiaoshoujia/models"
)

/*
投一票就加432分，86400/200, 那么天时间是86400，200张赞成票的帖子，才能在网站上置顶

direction = 1 时，有两种情况：
	1：之前没有投过票，现在投赞成票
	2：之前投反对票，现在改投赞成票

direction = 0 时，有两种情况：
	1：之前投过赞成票，现在取消投票
	2：之前投反对票，现在取消投票

direction = -1 时，有两种情况：
	1：之前没有投过票，现在投反对票
	2：之前投赞成票，现在改反对票

投票的限制：
	每个帖子一周以后就不能再次投票了
	到期之后，将redis中保存的赞成票数和反对票数存储到MySQL数据库
	到期之后，删除redis中keyPostVotedZSetPF这个key的分数
*/

// VoteForPost 为帖子投票的函数
func VoteForPost(userId int64, paramVoteData *models.ParamVoteData) {
	//zap.L().Debug("VoteForPost",
	//	zap.Int64("userID", userId),
	//	zap.String("postID", paramVoteData.PostId),
	//	zap.Int8("direction", paramVoteData.Direction))

	return redisconnect.VoteForPost()
}
