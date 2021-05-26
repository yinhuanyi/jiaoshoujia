/**
 * @Author: Robby
 * @File name: request.go
 * @Create date: 2021-05-22
 * @Function:
 **/

package controllers

import (
	"errors"
	"strconv"

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

// GetPagination  获取用户传递的分页数据
func GetPagination(c *gin.Context) (page, perPage int64, err error) {
	// 获取分页数据
	pageStr := c.Query("page")        // 第几页
	perPageStr := c.Query("per_page") // 返回几条数据

	// 这里如果用户没有彻底分页，那么类型换行会报错，那么直接使用默认值即可
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	perPage, err = strconv.ParseInt(perPageStr, 10, 64)
	if err != nil {
		perPage = 10
	}
	return
}
