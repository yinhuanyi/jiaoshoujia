/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package models

// 定义user表的表结构
type User struct {
	// 数据库字段的名称是user_id
	UserId   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
