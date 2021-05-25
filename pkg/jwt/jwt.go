/**
 * @Author: Robby
 * @File name: jwt.go
 * @Create date: 2021-05-20
 * @Function:
 **/

package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 过期时间
const TokenExpireDuration = time.Hour * 24

// 加盐
var mySecret = []byte("ipfs")

// jwt中的claim
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	// 这个是jwt的标准claim
	jwt.StandardClaims
}

// GenToken  生成JWT
func GenToken(userId int64, username string) (string, error) {

	// 创建一个我们自己的声明
	c := MyClaims{
		userId,
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "robby",                                    // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken ：解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 声明一个地址
	var mc = new(MyClaims)

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
