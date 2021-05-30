/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package models

import "time"

// Post 这个post结构体的字段，既需要与用户提交的参数对应，也需要与数据库中的表字段对应
// 这里将字段类型相同的放在一起，这种方式叫做【内存对齐】
type Post struct {
	// 如果需要将Post结构体，进行json序列化的时候，将int64类型转换为string类型，那么直接可以在tag后面加上string即可
	//Id          int64     `json:"id,string" db:"post_id"`
	Id          int64     `json:"id" db:"post_id"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityId int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail ： 定义返回的结构体数据
/*
最终返回的结构如下，post和community作为两个单独的字段，被json序列化了
{
    "code": 1000,
    "msg": "Success",
    "data": {
        "author_name": "yinhuanyi",
        "post": {
            "id": 119099726669811712,
            "author_id": 118699425387253760,
            "community_id": 2,
            "status": 0,
            "title": "日志报警",
            "content": "miner日志告警",
            "create_time": "2021-05-25T23:39:47Z"
        },
        "community": {
            "id": 2,
            "name": "leetcode",
            "introduction": "刷题刷题刷题",
            "create_time": "2020-01-01T08:00:00Z"
        }
    }
}
*/
type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	// 下面的Post和CommunityDetail指定了json结构化的方式，作为一个单独字段体现出来
	*Post            `json:"post"`
	*CommunityDetail `json:"community"`
}
