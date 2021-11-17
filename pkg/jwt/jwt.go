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

// jwt解析错误的类型常量
var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

// 过期时间
const TokenExpireDuration = time.Hour * 24

// 这是服务器端的secret
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

	// 创建一个返回secret的函数
	keyFunc := func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	}
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		// 判断token解析的错误类型
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 { // 不是一个jwt的token
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 { // token过期
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 { // token未激活
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid // token非法
			}
		}
		return nil, err
	}

	/*
		if token != nil {
			claims, ok := token.Claims.(*models.CustomClaims)
			if ok && token.Valid {
				return claims, nil
			}
			return nil, TokenInvalid
		} else {
			return nil, TokenInvalid
		}
	*/

	if token != nil {
		claims, ok := token.Claims.(*MyClaims)
		if ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}

	return nil, TokenInvalid

	// 校验token
	//if token.Valid {
	//	return mc, nil
	//}

	//return nil, errors.New("invalid token")
}
