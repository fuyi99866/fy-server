package core

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_server/conf"
	"go_server/docs"
	"go_server/models"
	"go_server/routers"
)

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

func RunWindowsServer()  {

	//读取配置文件
	config := flag.String("c", "./app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	conf.Init(*config) //根据配置文件初始化配置
	fmt.Println("服务器启动...")
	models.Init()

	if conf.Swag != nil {
		docs.SwaggerInfo.Host = conf.Swag.Host
		docs.SwaggerInfo.BasePath = conf.ServerSetting.BasePath
		scheme := "http"
		if conf.ServerSetting.HTTPS {
			scheme = "https"
		}
		logrus.Info(fmt.Sprintf("-----服务启动,可以打开  %s://%s%s/swagger/index.html 查看详细接口------",scheme,conf.Swag.Host,conf.ServerSetting.BasePath,))
	}

	routers.Run()

	//fmt.Println("-----服务启动------")
}