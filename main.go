package main

import (
	"go_server/core"
	alog "go_server/log"
)

func main() {
	//global.GVA_VP = core.Viper  //初始化Viper,加载配置到全局配置
	//global.GVA_LOG = core. //初始化zap日志库
    //global.GVA_DB = initialize.Gorm() // gorm连接数据库
	//初始化日志库
	alog.Init() //初始化日志库
	core.RunWindowsServer()


}
