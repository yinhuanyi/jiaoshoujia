/**
 * @Author: Robby
 * @File name: main.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package main

import (
	"context"
	"fmt"
	"jiaoshoujia/controllers"
	mysqlconnect "jiaoshoujia/dao/mysql"
	redisconnect "jiaoshoujia/dao/redis"
	"jiaoshoujia/pkg/snowflake"
	"jiaoshoujia/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"jiaoshoujia/logger"
	"jiaoshoujia/settings"
)

func main() {

	// 0：判断参数
	if len(os.Args) < 2 {
		fmt.Printf("运行一下 ./jiaoshoujia config.yaml")
		return
	}

	// 1：加载配置, 这里的os.Args[1] 就是传递的配置文件路径
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("配置文件读取失败：%v", err)
	}

	// 2：初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("日志配置失败：%v", err)
	}
	// 将缓冲区的日志写入日志文件
	defer zap.L().Sync()

	// 3：初始化MySQL数据库
	if err := mysqlconnect.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("MySQL连接失败：%v", err)
	}
	defer mysqlconnect.Close() // 程序停止后，关闭连接

	// 4：初始化Redis
	if err := redisconnect.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("Redis连接失败：%v", err)
	}
	defer redisconnect.Close() // 程序停止后，关闭连接

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId); err != nil {
		fmt.Printf("初始化雪花算法失败：%v", err)
	}

	// 5：注册路由
	r := routes.Init(settings.Conf.Mode)

	// 6：初始化gin框架校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("翻译器获取失败：%v", err)
	}

	// 7：启动服务 (下面是优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 虽然这里设置了5秒超时，但是手动调用，
	defer cancel()
	// srv.Shutdown(ctx)中ctx.Done()是阻塞的
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
