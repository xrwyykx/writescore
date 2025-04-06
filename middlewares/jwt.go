package middlewares

import (
	"crypto/rand"
	"github.com/golang-jwt/jwt"
	"io"
	"time"
)

// 随机生成一个密匙
func generateJWTSecret() ([]byte, error) {
	const keyLength = 32
	key := make([]byte, keyLength)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 自定义Claims结构体,用于存储kwt的自定义声明，
type CustomClaims struct {
	UserId             int64 `json:"user_id"`
	jwt.StandardClaims       //用于存储标准的kwt声明，包含令牌的过期时间等
}

var jwtsecret []byte //用于存储生成的jwt密匙

// 生成jwt token
func GEnerateJWT(userId int64) (string, error) {
	var err error
	jwtsecret, err = generateJWTSecret() //生成密匙
	if err != nil {
		return "", err
	}
	claims := CustomClaims{UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), //有效期为一个小时
		},
	} //创建claims实例
	//使用HS256签名方法和密匙生成jwt令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(jwtsecret)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

// 解析jwt
func PraseToken(tokenstring string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	//接收一个令牌字符串，解析令牌的有效性
	token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtsecret, nil
	})
	//出现错误或者令牌无效就返回错误
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
