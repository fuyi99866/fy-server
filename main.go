package main

import (
	"flag"
	"fmt"
	"go_server/conf"
	"go_server/core"
	"go_server/docs"
	"go_server/logger"
	"go_server/models"
)

func main() {

	//读取配置文件
	config := flag.String("c", "./app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	conf.Init(*config) //根据配置文件初始化配置

	//初始化日志系统
	//alog.Init() //初始化日志库 ,使用logrus库
	logger.InitLog1(conf.AppSetting.LogLever,"/logs/go_server.log") //初始化日志库 ,使用zap库

	//初始化数据库
	models.Init()

	if conf.Swag != nil {
		docs.SwaggerInfo.Host = conf.Swag.Host
		docs.SwaggerInfo.BasePath = conf.ServerSetting.BasePath
		scheme := "http"
		if conf.ServerSetting.HTTPS {
			scheme = "https"
		}
		logger.Info(fmt.Sprintf("-----服务启动,可以打开  %s://%s%s/swagger/index.html 查看详细接口------",scheme,conf.Swag.Host,conf.ServerSetting.BasePath,))
	}

	core.RunWindowsServer()
}
