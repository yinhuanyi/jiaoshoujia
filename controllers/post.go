/**
 * @Author: Robby
 * @File name: post.go
 * @Create date: 2021-05-25
 * @Function:
 **/

package controllers

import (
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

// GetPostDetail : 获取帖子详情
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
