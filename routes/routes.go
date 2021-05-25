/**
 * @Author: Robby
 * @File name: routes.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package routes

import (
	"jiaoshoujia/controllers"
	"jiaoshoujia/logger"
	"jiaoshoujia/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(mode string) *gin.Engine {
	// 如果配置文件中，设置为release模式，那么gin就是release
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 在路由层将日志中间件加入
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 路由分组
	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controllers.SignUpHandler)
	// 在用户登录的时候返回jwt token
	v1.POST("/login", controllers.LoginHandler)
	// 应用jwt中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	// Todo: 为什么要加上这花括号
	{
		// 获取列表
		v1.GET("/community", controllers.CommunityHandler)
		// 获取详情
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		// 创建一个帖子
		v1.POST("/post", controllers.CreatePostHandler)
		// 获取帖子详情
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
	}

	// JWTAuthMiddleware函数，对jwt进行校验(anthorization头部进行验证)
	// 用户请求头携带标准的Bearer Auth访问/ping接口
	r.GET("/ping", func(context *gin.Context) {
		isLogin := true
		if isLogin {
			context.String(http.StatusOK, "pong")
		} else {
			context.String(http.StatusOK, "请登录")
		}
	})

	return r
}
