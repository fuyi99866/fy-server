package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"net/http"
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
	logger.Info("GenerateToken ")
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

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		Authorization := context.GetHeader("Authorization")//验证token，要从Header中查询Authorization
		token := Authorization
		fmt.Println("jwt", token)
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := ParseToken(token)
			fmt.Println("解析出来的claims:", claims)
			fmt.Println("解析出来的err:", err)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code!= e.SUCCESS {
			context.JSON(http.StatusUnauthorized,gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
