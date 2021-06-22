package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var engine *gin.Engine

func run()  {
	//1.创建路由
	r := gin.Default()
	//2、绑定路由规则，执行的函数
	//gin.Context,封装了request和response
	r.GET("/user/:name/*action", func(context *gin.Context) {
		name := context.Param("name")
		action := context.Param("action")
		//截取
		action = strings.Trim(action,"/")

		context.String(http.StatusOK,name+" is "+action)
	})
	//3、监听端口，默认是8080
	//Run("里面不指定端口就默认为8080")
	r.Run(":8080")
}
