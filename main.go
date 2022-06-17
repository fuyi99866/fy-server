package main

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/sirupsen/logrus"
	"go_server/conf"
	"go_server/cron"
	"go_server/docs"
	"go_server/models"
	"go_server/pkg/gredis"
	"go_server/pkg/log"
	"go_server/pkg/setting"
	"go_server/routers"
	"go_server/service"
	"syscall"
	"time"
)

/**
  执行 swag init --generalInfo .\main.go 生成docs
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

// @title code server
// @version 0.0.1
// @description Go 学习综合demo
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	//读取配置文件
	config := flag.String("c", "./conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()//把用户传递的命令行参数解析为对应的变量值
	setting.Init(*config) //根据配置文件初始化配置

	//初始化日志系统
	log.Init()

	//初始化数据库
	models.Init()
	models.Casbin_Init()

	gredis.InitRedis()

	//TODO 启动MQTT服务
	go service.Start()

	initServer()

	//开始定时任务
	go cron.Start()

}


//初始化服务
func initServer() {
	//注册路由
	app := routers.InitRouter()

	//热更新
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout*time.Second
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout*time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, app)
	server.BeforeBegin = func(add string) {
		conf.Pid = syscall.Getpid()
		logrus.Info("Actual pid is ", syscall.Getpid())
	}

	if setting.Swag != nil {
		docs.SwaggerInfo.Host = setting.Swag.Host
		docs.SwaggerInfo.BasePath = setting.ServerSetting.BasePath
		scheme := "http"
		if setting.ServerSetting.HTTPS {
			scheme = "https"
		}
		logrus.Info(fmt.Sprintf("-----服务启动,可以打开  %s://%s%s/swagger/index.html 查看详细接口------", scheme, setting.Swag.Host, setting.ServerSetting.BasePath, ))
	}

	if setting.ServerSetting.HTTPS {
		//支持https
		server.ListenAndServeTLS("","")
	}else {
		server.ListenAndServe()
	}
}


