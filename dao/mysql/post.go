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

// GetPostList : 获取post列表数据
func GetPostList(page, perPage int64) (postList []*models.Post, err error) {

	// postList = new([]*models.Post) 这样声明切片指针是错误的

	// 初始化列表，使用make函数，给一个cap为2即可
	postList = make([]*models.Post, 0, 2) // 使用make初始化的时候，不要写成make([]*models.Post, 2)，这样会多了个nil元素
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post order by create_time desc limit ?,?"
	err = db.Select(&postList, sqlStr, (page-1)*perPage, perPage)
	if err != nil {
		zap.L().Error("db.Select(&postList, sqlStr)", zap.Error(err))
	}
	return
}
