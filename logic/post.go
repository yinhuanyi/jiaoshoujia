/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package logic

import (
	mysqlconnect "jiaoshoujia/dao/mysql"
	"jiaoshoujia/models"
	"jiaoshoujia/pkg/snowflake"
)

func CreatePost(post *models.Post) (err error) {
	// 1：生成post id, 赋值给post.id，由于这里是地址，因此post的值会保存下来
	post.Id = snowflake.GenID()
	// 2：将post数据写入数据库
	mysqlconnect.CreatePost(post)
	return
}
