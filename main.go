package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_server/conf"
	"go_server/core"
	"go_server/docs"
	alog "go_server/log"
	"go_server/models"
)

func main() {
	//global.GVA_VP = core.Viper  //初始化Viper,加载配置到全局配置
	//global.GVA_LOG = core. //初始化zap日志库
    //global.GVA_DB = initialize.Gorm() // gorm连接数据库

	//读取配置文件
	config := flag.String("c", "./app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	conf.Init(*config) //根据配置文件初始化配置

	//初始化日志系统
	//alog.Init() //初始化日志库 ,使用logrus库
	alog.InitLog1()//初始化日志库 ,使用zap库

	//初始化数据库
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

	core.RunWindowsServer()
}
