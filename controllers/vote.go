/**
 * @Author: Robby
 * @File name: vote.go
 * @Create date: 2021-05-27
 * @Function:
 **/

package controllers

import (
	"jiaoshoujia/logic"
	"jiaoshoujia/models"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PostVoteController 用户投票请求
func PostVoteController(c *gin.Context) {
	// 返回结构体类型指针
	paramVoteData := new(models.ParamVoteData)
	err := c.ShouldBindJSON(paramVoteData)
	// 下面是对错误的返回处理
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(paramVoteData)", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 判断是不是json解析错误
		if !ok {                                     // 如果是json解析错误
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	// 基于context获取用户的ID
	userId, err := GetCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 将用户的ID和用户请求的参数传递到VoteForPost函数
	if err = logic.VoteForPost(userId, paramVoteData); err != nil {
		zap.L().Error("logic.VoteForPost(userId, paramVoteData)", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
		//ResponseError(c, CodeServerBusy)
		//return
	}

	ResponseSuccess(c, nil)
}
