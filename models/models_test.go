package models

import (
	"flag"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"testing"
)

type Result struct {
	Article
	TagName  string `json:"tag_name"`
	TagState string `json:"tag_state"`
}




func TestAll(t *testing.T) {
	//读取配置文件
	config := flag.String("c", "D:/fy/git/go_server/fy-server/conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置
	logger.InitLog1(setting.AppSetting.LogLever, "./logs/go_server.log") //初始化日志库 ,使用zap库
	//初始化数据库
	Init()
}



