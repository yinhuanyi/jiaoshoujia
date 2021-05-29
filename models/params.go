/**
 * @Author: Robby
 * @File name: params.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

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
	PostId string `json:"post_id" binding:"required"` // 帖子ID
	// 这个字段有两个约束条件，必须是 -1 0 1 中的一个, 这个字段不用required，因为传递0的时候，会被任务没有传递数据，默认int8的值就是0
	Direction int8 `json:"direction,string" binding:"oneof=-1 0 1"` // 赞成票(1)或反对票(-1)取消投票(0)
}

// ParamPostList :帖子列表详情请求参数
type ParamPostList struct {
	// form表示的是解析form类型的请求参数，也就是url参数
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}
