/**
 * @Author: Robby
 * @File name: keys.go
 * @Create date: 2021-05-27
 * @Function:
 **/

package redisconnect

// redis key 使用命名空间的方式进行区分，使用:作为名称空间的分隔符
const (
	Prefix             = "ipfsmain:"   // 前缀
	KeyPostTimeZSet    = "post:time"   // zset类型，帖子及发帖时间
	keyPostScoreZSet   = "post:score"  // zset类型，帖子及投票的分数
	keyPostVotedZSetPF = "post:voted:" // zset类型，记录用户及投票类型， 参数是post id, 这个key在存储的时候会加上postId
	KeyCommunitySetPF  = "community:"  // set类型;保存每个分区下帖子的ID
)

// 获取带前缀的key
func getRedisKey(key string) (newkey string) {
	return Prefix + key
}
