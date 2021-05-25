/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package logic

import (
	mysqlconnect "jiaoshoujia/dao/mysql"
	"jiaoshoujia/models"
	"jiaoshoujia/pkg/jwt"
	"jiaoshoujia/pkg/snowflake"
)

// SignUp 用户注册
func SignUp(params *models.ParamSignUp) (err error) {

	// 判断用户是否存在
	err = mysqlconnect.CheckUserExist(params.Username)
	if err != nil {
		// 查库错误
		return
	}

	// 1：生成UID
	userId := snowflake.GenID()

	// 2：创建user表记录
	userInstance := models.User{
		UserId:   userId,
		Username: params.Username,
		Password: params.Password,
	}

	// 3：记录保存到数据库
	err = mysqlconnect.InsertUser(&userInstance)
	return
}

// Login 用户登录
func Login(params *models.ParamLogin) (token string, err error) {
	// 创建记录
	user := models.User{
		Username: params.Username,
		Password: params.Password,
	}
	// 由于这里传递的是指针，因此，可以基于user变量获取到用户的userId
	err = mysqlconnect.Login(&user)
	// 这里相当于return err
	// 如果登录失败，返回
	if err != nil {
		return
	}
	// 返回一个token和err
	return jwt.GenToken(user.UserId, user.Username)
}
