package core

import (
	"github.com/sirupsen/logrus"
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

	logrus.Debug("--------服务启动-------")

	routers.Run()

}