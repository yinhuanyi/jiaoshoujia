/**
 * @Author: Robby
 * @File name: params.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package models

// ParamSignUp 【请求进来】定义请求的参数结构体, 注册请求参数
type ParamSignUp struct {
	// 这个binding tag是gin中validator中的，用于对参数进行校验
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData :投票数据
type ParamVoteData struct {
	PostId int64 `json:"post_id,string" binding:"required"` // 帖子ID
	// 这个字段有两个约束条件，必须是 -1 0 1 中的一个
	Direction int8 `json:"direction,string" binding:"required,oneof=-1 0 1"` // 赞成票(1)或反对票(-1)取消投票(0)
}
