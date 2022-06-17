package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go_server/pkg/setting"
	"net/http"
	"strings"
	"time"
)

//打印路由
func Log2File() gin.HandlerFunc {
	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUri := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求IP
		clientIP := c.ClientIP()
		if reqUri == (setting.ServerSetting.BasePath+"/alive") ||
			reqUri == (setting.ServerSetting.BasePath+"/swagger/index.html") ||
			reqUri == (setting.ServerSetting.BasePath+"/version") {
			//非必要信息不打印
		} else {
			// 日志格式
			logrus.Infoln(fmt.Sprintf("| %3d | %13v | %15s | %s | %s ",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri),
			)
		}
	}
}

//允许跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method //请求方法
		var headerKeys []string    //声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Origin", "*")                                   //也是允许访问所有域
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE") //服务器支持所有跨域请求的方法，为了避免浏览请求的多次“预检”请求
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept,"+
			" Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive,"+
			" User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma") //header的类型 可以返回其他字段
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
			"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") //跨域关键设置 让浏览器可以解析
		c.Header("Access-Control-Max-Age", "172800")                                                                        //缓存请求信息 单位为秒
		c.Header("Access-Control-Allow-Credentials", "false")                                                               //跨域请求是否需要带cookie 默认设置为true
		c.Set("content-type", "application/json")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next() //处理请求
	}
}

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:9096",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		// 继续往下处理
		c.Next()
	}
}
