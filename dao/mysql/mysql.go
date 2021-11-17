/**
 * @Author: Robby
 * @File name: mysql.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package mysqlconnect

import (
	"fmt"
	"jiaoshoujia/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	//db, err = sqlx.Connect("mysql", dsn)
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("连接数据库失败: %v", r)
			panic(r)
		}
	}()
	// 如果连接不上，这里会panic
	db = sqlx.MustConnect("mysql", dsn)

	//if err != nil {
	//	fmt.Printf("connect DB failed, err:%v\n", err)
	//	zap.L().Error("connect DB failed", zap.Error(err))
	//	return
	//}
	// SetMaxIdleConns 设置空闲的最大连接数连接池。
	//
	// 如果 MaxOpenConns 大于 0 但小于新的 MaxIdleConns，
	// 然后将减少新的 MaxIdleConns 以匹配 MaxOpenConns 限制。
	//
	// 如果 n <= 0，则不保留空闲连接。
	//
	// 默认的最大空闲连接数当前为 2。这可能会改变
	// 未来版本。
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return

}

func Close() {
	_ = db.Close()
}
