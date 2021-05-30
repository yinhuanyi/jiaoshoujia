/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package mysqlconnect

import (
	"jiaoshoujia/models"
	"strings"

	"github.com/jmoiron/sqlx"

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
	// 这里一定要再次获取一下post的结构体类型指针，要不然会报错：【nil pointer passed to StructScan destination】
	// 这里的参数post是一个空指针，因此必须要重新赋值，使用new函数创建一个结构体的类型指针
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
	//postList = make([]*models.Post, 0, 2) // 使用make初始化的时候，不要写成make([]*models.Post, 2)，这样会多了个nil元素
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post order by create_time desc limit ?,?"
	err = db.Select(&postList, sqlStr, (page-1)*perPage, perPage)
	if err != nil {
		zap.L().Error("db.Select(&postList, sqlStr)", zap.Error(err))
	}
	return
}

// GetPostListByIds ：根据给定的ID列表查询帖子的数据
func GetPostListByIds(ids []string) (postList []*models.Post, err error) {
	// FIND_IN_SET函数在SQL中可以让post_id字段，按照ids的顺序排序
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
`
	// 这里可以认为是格式化, args返回的是一个列表
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// 需要重新绑定一下
	query = db.Rebind(query)
	// postList = make([]*models.Post, 0, 2)  这里不用make初始化，因为切片声明后，再取地址，不是空指针
	err = db.Select(&postList, query, args...)

	return
}
