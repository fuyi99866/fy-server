package jwt

import (
	"github.com/gin-gonic/gin"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/util"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		Authorization := context.GetHeader("Authorization")//验证token，要从Header中查询Authorization
		token := Authorization
		logger.Info("jwt", token)
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			logger.Info("解析出来的claims:", claims)
			logger.Info("解析出来的err:", err)
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
