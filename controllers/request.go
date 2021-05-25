/**
 * @Author: Robby
 * @File name: request.go
 * @Create date: 2021-05-22
 * @Function:
 **/

package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const ContextUserIdKey = "userId"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUserId ：获取当前请求用户的userid
func GetCurrentUserId(c *gin.Context) (userId int64, err error) {
	uid, ok := c.Get(ContextUserIdKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	// 这里是使用了断言，断言也会有错误，因此有一个OK判断是否断言成功
	userId, ok = uid.(int64)

	if !ok {
		err = ErrorUserNotLogin
		return
	}

	return
}
