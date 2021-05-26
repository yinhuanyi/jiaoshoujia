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

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 1：生成post id, 赋值给post.id，由于这里是地址，因此post的值会保存下来
	post.Id = snowflake.GenID()
	// 2：将post数据写入数据库
	err = mysqlconnect.CreatePost(post)
	return
}

func GetPostById(postId int64) (data *models.ApiPostDetail, err error) {
	//data = new(models.ApiPostDetail)
	// 基于postid查询post详情，并且基于详情组装 ApiPostDetail类型的数据返回
	postInstance, err := mysqlconnect.GetPostById(postId)
	if err != nil {
		zap.L().Error("mysqlconnect.GetPostById(postId)", zap.Int64("postId", postId), zap.Error(err))
		return
	}

	// 基于userid获取user相关信息
	userInstance, err := mysqlconnect.GetUserById(postInstance.AuthorId)
	if err != nil {
		zap.L().Error("mysqlconnect.GetUserById(postInstance.AuthorId)", zap.Int64("AuthorId", postInstance.AuthorId), zap.Error(err))
		return
	}

	// 基于comunityId获取community信息
	communityInstance, err := mysqlconnect.GetCommunityDetailById(postInstance.CommunityId)
	if err != nil {
		zap.L().Error("mysqlconnect.GetCommunityDetailById(postInstance.CommunityId)", zap.Int64("CommunityId", postInstance.CommunityId), zap.Error(err))
		return
	}

	// 拼凑最终返回的数据, 由于data最终会被序列化为json，因此在定义ApiPostDetail结构体时候，可以指定结构体的json化方式
	data = &models.ApiPostDetail{
		AuthorName:      userInstance.Username,
		Post:            postInstance,
		CommunityDetail: communityInstance,
	}

	return

}

func GetPostList(page, perPage int64) (data []*models.ApiPostDetail, err error) {

	// 查询数据库, 获取列表, 这里不是直接返回数据，还需要对数据进行聚合
	postList, err := mysqlconnect.GetPostList(page, perPage)
	if err != nil {
		return
	}

	// 对data列表数据，进行初始化, len(postList): 表示返回了多少帖子数，那么就初始化多少
	data = make([]*models.ApiPostDetail, 0, len(postList))

	// 由于需要返回多个ApiPostDetail实例，那么使用for循环遍历
	for _, postInstance := range postList {
		// 基于userid获取user相关信息
		userInstance, err := mysqlconnect.GetUserById(postInstance.AuthorId)
		if err != nil {
			zap.L().Error("mysqlconnect.GetUserById(postInstance.AuthorId)", zap.Int64("AuthorId", postInstance.AuthorId), zap.Error(err))
			// 这里在循环里面，不能写return，如果有错误返回，那么直接continue，进行下一轮循环
			continue
		}

		// 基于comunityId获取community信息
		communityInstance, err := mysqlconnect.GetCommunityDetailById(postInstance.CommunityId)
		if err != nil {
			zap.L().Error("mysqlconnect.GetCommunityDetailById(postInstance.CommunityId)", zap.Int64("CommunityId", postInstance.CommunityId), zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      userInstance.Username,
			Post:            postInstance,
			CommunityDetail: communityInstance,
		}

		data = append(data, postDetail)

	}

	return

}
