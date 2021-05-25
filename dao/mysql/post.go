/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package mysqlconnect

import (
	"jiaoshoujia/models"

	"go.uber.org/zap"
)

// CreatePost 将post结构体数据写入到数据表
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.Id, post.Title, post.Content, post.AuthorId, post.CommunityId)
	if err != nil {
		zap.L().Error("insert into post(post_id, title, content, author_id, community_id) values (?,?,?,?,?)", zap.Error(err))
	}
	return
}

// GetPostById 获取post详情
func GetPostById(postId int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post where post_id = ?"
	err = db.Get(post, sqlStr, postId)
	if err != nil {
		zap.L().Error("db.Get(post, sqlStr, postId)", zap.Error(err))
	}
	return
}
