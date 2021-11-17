/**
 * @Author: Robby
 * @File name: logger.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package logger

import (
	"jiaoshoujia/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

func Init(cfg *settings.LogConfig, mode string) (err error) {

	//writeSyncer := getLogWriter(viper.GetString("log.filename"), viper.GetInt("log.file_size"), viper.GetInt("log.max_backups"), viper.GetInt("max_size"))
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}

	var core zapcore.Core
	// 如果项目在开发模式下, 让日志既可写入到文件，也可以输出到终端
	if mode == "dev" {
		// 日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg = zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L().debug()调用即可
	zap.ReplaceGlobals(lg)
	return

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 下面两个日志配置，会在路由中间件中被用到

// GinLogger 接收gin框架默认的日志, 这个也是中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 直接调用第二个中间件
		c.Next()

		// 当请求处理完毕后，再执行下面的语句，因此这里可以计算出请求消耗的时间
		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover项目可能出现的panic，并使用zap记录相关日志, 这个GinRecovery函数是一个中间件
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用recover的时候，defer的函数是自执行函数，这里是所有请求最后执行这个defer后的函数的，保证可以捕获到panic异常
		defer func() {

			if err := recover(); err != nil {
				// 这里是检查是不是brokenpipe的错误
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取客户端的请求路径
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 如果是brokenpipe问题
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error))
					// 不再进入下一个中间件，当前中间件运行完毕即可
					c.Abort()
					// return 下面的语句也不会执行了
					return
				}

				// 记录堆栈信息
				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 返回500响应报文
				c.AbortWithStatus(http.StatusInternalServerError)
			}

		}()

		c.Next()
	}
}
