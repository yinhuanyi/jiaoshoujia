/**
 * @Author: Robby
 * @File name: response.go
 * @Create date: 2021-05-19
 * @Function:  这里是封装一些响应状态码
 **/

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	{
		"code": 0001, 业务状态码
		"msg": "提示信息",
		"data": {} 数据
	}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 返回错误信息，带简单的错误提示
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.GetMsg(),
		Data: nil,
	}

	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 返回错误信息，有自定义的msg提示
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	c.JSON(http.StatusOK, rd)
}

// ResponseSuccess 返回成功的信息，带data数据
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.GetMsg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
