package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go_server/pkg/setting"
	"time"
)

var JwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//产生token的函数
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	//设置失效时间
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	//指明生成算法，生成token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	logrus.Info("GenerateToken ")
	token, err := tokenClaims.SignedString(JwtSecret)

	return token, err
}

//验证token的函数
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}


