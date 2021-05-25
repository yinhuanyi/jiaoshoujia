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

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler : 创建帖子
func CreatePostHandler(c *gin.Context) {

	// 0：申明结构体指针
	post := new(models.Post)

	// 1：获取参数，参数校验
	if err := c.ShouldBindJSON(post); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2：创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post)", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}

	// 3：返回响应
	ResponseSuccess(c, nil)

}
