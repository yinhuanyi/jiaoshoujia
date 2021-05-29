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
	"strconv"

	"go.uber.org/zap"
)

// VoteForPost 为帖子投票的函数
func VoteForPost(userId int64, paramVoteData *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userId),
		zap.String("postID", paramVoteData.PostId),
		zap.Int8("direction", paramVoteData.Direction))

	return redisconnect.VoteForPost(strconv.Itoa(int(userId)), paramVoteData.PostId, float64(paramVoteData.Direction))
}
