/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-05-19
 * @Function:
 **/

package mysqlconnect

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"jiaoshoujia/models"
)

const secret = "ipfsmain"

// CheckUserExist 判断用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	// 这里的count是变量, 一般情况下，如果是查询一个记录，一般是一个结构体
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser (err error)相当于申明了var err error，所有return return的是err
func InsertUser(userInstance *models.User) (err error) {
	// 对密码进行加密
	userInstance.Password = encryptPassword(userInstance.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values (?,?,?)`
	_, err = db.Exec(sqlStr, userInstance.UserId, userInstance.Username, userInstance.Password)
	return
}

// encryptPassword 对密码进行加密
func encryptPassword(password string) string {
	h := md5.New()
	// 加盐
	h.Write([]byte(secret))
	// 加密
	return hex.EncodeToString(h.Sum([]byte(password)))
}

// Login 判断用户传递的密码，是否与密码相等
func Login(user *models.User) (err error) {

	opassword := user.Password // 用户传递进来的password
	sqlStr := `select user_id, username, password from user where username = ?`
	// 这里其实会覆盖之前的密码，因为查询的结构赋值给了user这个结构体
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows { // 如果没有查到用户
		return ErrorUserNotExist
	}
	if err != nil { // 如果是其他错误
		return
	}
	// 判断密码是否正确
	password := encryptPassword(opassword) // 获取加密后的密码
	if password != user.Password {         // 判断加密后的密码，是否与数据库中的密码相等
		return ErrorInvalidPassword
	}
	// 如果什么错误都没有，那么return相当于return err，此时err是nil
	return

}
