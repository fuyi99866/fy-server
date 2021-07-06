package main

import (
	"flag"
	"fmt"
	"go_server/docs"
	"go_server/models"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"go_server/routers"
	"net/http"
	"time"
)

/**
  执行 swag init --generalInfo .\routers\http.go 生成docs
  https://www.ctolib.com/swaggo-swag.html
*/

/**
1、数据库操作gorm
2、服务器接口框架gin
3、MQTT
4、HTTP和websocket
5、日志打印logrus
6、路由、校验Token jwt
7、生成标准的在线文档swagger
8、casbin控制访问权限
9、使用单元测试进行代码性能检测
*/
func main() {

	//读取配置文件
	config := flag.String("c", "./app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置

	//初始化日志系统
	//alog.Init() //初始化日志库 ,使用logrus库
	logger.InitLog1(setting.AppSetting.LogLever, "/logs/go_server.log") //初始化日志库 ,使用zap库

	//初始化数据库
	models.Init()

	initServer()
}

//初始化服务
func initServer() {
	app := routers.InitRouter()

	//注册路由
	routers.RegisterRouter(app)

	initHTTPServer(app)
}

// InitHTTPServer 初始化http服务
func initHTTPServer(handler http.Handler) {
	logger.Info("start server")
	if setting.Swag != nil {
		docs.SwaggerInfo.Host = setting.Swag.Host
		docs.SwaggerInfo.BasePath = setting.ServerSetting.BasePath
		scheme := "http"
		if setting.ServerSetting.HTTPS {
			scheme = "https"
		}
		logger.Info(fmt.Sprintf("-----服务启动,可以打开  %s://%s%s/swagger/index.html 查看详细接口------", scheme, setting.Swag.Host, setting.ServerSetting.BasePath, ))
	}

	srv := &http.Server{
		Addr:         ":" + fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
