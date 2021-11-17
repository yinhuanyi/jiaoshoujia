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
	"time"

	"github.com/gin-gonic/gin"
)

func Init(mode string) *gin.Engine {
	// 如果配置文件中，设置为release模式，那么gin就是release
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	/*
		这个默认会加入Logger(), Recovery()两个中间件
		Logger() 中间件会请求的时候，将日志显示到终端
		Recovery() 中间件，当系统崩溃，会显示500的状态码
	*/
	//r := gin.Default()
	r := gin.New()

	// 在路由层将日志中间件加入
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 路由分组
	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controllers.SignUpHandler)
	// 在用户登录的时候返回jwt token
	v1.POST("/login", controllers.LoginHandler)

	// 应用jwt中间件, 限流中间件
	v1.Use(middlewares.JWTAuthMiddleware(), middlewares.RateLimitMiddleware(2*time.Second, 1))

	{
		// 获取列表
		v1.GET("/community", controllers.CommunityHandler)
		// 获取详情
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		// 创建一个帖子
		v1.POST("/post", controllers.CreatePostHandler)
		// 获取帖子详情
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		// 帖子的列表(默认使用MySQL中帖子创建时间排序)
		v1.GET("/posts", controllers.GetPostListHandler)
		// 帖子的列表(使用Redis中帖子创建时间或帖子分数排序)
		v1.GET("/posts2", controllers.GetPostListHandler2)
		// 投票
		v1.POST("/vote", controllers.PostVoteController)
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
