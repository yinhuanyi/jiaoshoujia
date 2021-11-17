/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package controllers

import (
	"fmt"
	"jiaoshoujia/logic"
	"jiaoshoujia/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler : 创建帖子
func CreatePostHandler(c *gin.Context) {

	// 0：申明结构体指针
	post := new(models.Post)

	// 1：获取参数，参数校验
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Error("c.ShouldBindJSON(post)", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 如果用户能够post帖子，那么一定是登录过的，因为路由的时候会经过jwt解析，那么用户的userid会写入到context里面
	userId, err := GetCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 将用户的ID赋值给post结构体
	post.AuthorId = userId

	// 2：创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post)", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}

	// 3：返回响应
	ResponseSuccess(c, nil)

}

// GetPostDetailHandler : 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {

	// 1：获取帖子的ID
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(postIdStr, 10, 64)", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2：根据帖子ID查询数据
	data, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPostById(postId)", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3：返回数据
	ResponseSuccess(c, data)

}

// GetPostListHandler : 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 这里将用户的分页数据提取出去了
	page, perPage, _ := GetPagination(c)
	// 获取数据， 将分页数据传递进来
	data, err := logic.GetPostList(page, perPage)
	if err != nil {
		zap.L().Error("logic.GetPostList()", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 : 获取帖子详情列表，升级版本
func GetPostListHandler2(c *gin.Context) {
	// c.ShouldBind() : 根据请求数据类型content-Type自动判断
	// c.ShouldBindJSON() : 解析json相关数据

	// 1：获取请求参数
	p := &models.ParamPostList{ // 这种方式可以设置默认参数, 之前都是申明一个结构体类型指针 p := new(models.ParamPostList)
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// 这里的ShouldBindQuery是获取参数的
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("c.ShouldBindQuery(p)", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2：去Redis查询ID列表, 根据ID去数据库查询帖子详情信息
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2(p)", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// GetCommunityPostListHandler 请求中，如果可以复用结构体，那么可以像这样复用
func GetCommunityPostListHandler(c *gin.Context) {

	p := &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
		CommunityId: 1,
	}
	fmt.Println(p)
}
