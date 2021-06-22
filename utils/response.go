package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int   //错误码
	Msg string //错误信息
	Data string `json:"data"` //详细错误信息
}

//简单的回复成功
func Response_test(context *gin.Context) {
	message := "成功"
	code := 200
	data := "data"
	context.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
		"result":  "true",
	})
}

func (context *Gin)Response(httpCode,errCode int,data interface{})  {
	context.C.JSON(httpCode,gin.H{
		"code":errCode,
		"msg":GetMsg(errCode),
		"data":data,
	})
	return
}
