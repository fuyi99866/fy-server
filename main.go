package main

import (
	"fmt"
	"go_server/core"
)

func main() {
	//global.GVA_VP = core.Viper  //初始化Viper,加载配置到全局配置
	//global.GVA_LOG = core.Zap() //初始化zap日志库
    //global.GVA_DB = initialize.Gorm() // gorm连接数据库
	fmt.Println("start")
	core.RunWindowsServer()
}
