package app

import (
	"github.com/gin-gonic/gin"
	"go_server/pkg/e"
)


type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int   //错误码
	Msg string //错误信息
	Data string `json:"data"` //详细错误信息
}


func (context *Gin)Response(httpCode,errCode int,data interface{})  {
	context.C.JSON(httpCode,gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
	return
}

func (g *Gin) ResponseTest(httpCode int, data interface{},errCode int,) {
	g.C.JSON(httpCode, gin.H{
		"respCode": 10000000,
		"respMessage": "操作成功！",
		"data": data,
	})

	return
}


