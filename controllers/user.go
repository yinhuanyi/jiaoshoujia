/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package controllers

import (
	"errors"
	mysqlconnect "jiaoshoujia/dao/mysql"
	"jiaoshoujia/logic"
	"jiaoshoujia/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 注册
func SignUpHandler(c *gin.Context) {
	// 1：校验参数
	//var params models.ParamSignUp
	params := new(models.ParamSignUp)

	if err := c.ShouldBindJSON(params); err != nil {
		zap.L().Error("注册方法，参数校验错误", zap.Error(err))
		// 判断是不是validator.ValidationErrors, 如果不是，那么不需要翻译，直接返回即可
		// 如果错误类型是ValidationErrors，错误将会被转换为errs错误，使用validator库对请求参数进行校验
		// 这里可能需要搞清楚validator校验：https://www.liwenzhou.com/posts/Go/validator_usages/
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 如果是validator.ValidationErrors类型错误，就可以直接翻译错误即可
		//c.JSON(http.StatusOK, gin.H{
		//	// 将翻译的错误返回
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		// 这里return是避免代码再向下执行
		return
	}

	// 2：对业务字段进行校验：判断参数是否为空，是否密码等于确认密码
	//if len(params.Username) == 0 || len(params.Password) == 0 || len(params.RePassword) == 0 || params.Password != params.RePassword {
	//	zap.L().Error("请求参数错误", zap.Any("params", params))
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	// 这里return是避免代码再向下执行
	//	return
	//}

	// 2：将用户信息写入数据库
	if err := logic.SignUp(params); err != nil {
		zap.L().Error("注册错误", zap.Error(err))
		if errors.Is(err, mysqlconnect.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3：返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "注册成功",
	//})
	ResponseSuccess(c, nil)
}

// 登录
func LoginHandler(c *gin.Context) {
	// 1：获取请求参数及参数校验
	params := new(models.ParamLogin)
	if err := c.ShouldBindJSON(params); err != nil {
		zap.L().Error("登录参数校验错误", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam) // 使用自定义错误返回
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 使用自定义错误返回带自定义的msg
		return
	}

	// 2：业务逻辑处理
	token, err := logic.Login(params)
	if err != nil {
		zap.L().Error("登录失败", zap.String("username", params.Username), zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "登录失败, 用户名或密码错误",
		//})
		// 下面是判断错误是什么类型的
		if errors.Is(err, mysqlconnect.ErrorUserNotExist) { // 如果返回的错误是用户不存在错误
			ResponseError(c, CodeUserNotExist)
		} else if errors.Is(err, mysqlconnect.ErrorInvalidPassword) { // 如果是密码错误
			ResponseError(c, CodeInvalidPassword)
		} else {
			ResponseError(c, CodeServerBusy)
		}

		return
	}

	// 3：返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功",
	//})
	ResponseSuccess(c, token)
}
