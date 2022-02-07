package adis_server

import (
	"flag"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"testing"
)

func TestAll(t *testing.T) {
	logger.InitLog(setting.AppSetting.LogLever, "./logs/log.txt") //初始化日志库 ,使用zap库
	config := flag.String("c", "D:/fy/git/go_server/fy-server/conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置

	//初始化队列，创建生产者和消费者
	logger.Info("setting.RocketMq.Addr: ", setting.RocketMq.Addr)
	logger.Info("setting.RocketMq.GroupName: ", setting.RocketMq.GroupName)
	mq := NewMq(WithAddr(setting.RocketMq.Addr),
		WithGroupName(setting.RocketMq.GroupName),
		WithRetry(2),
	)
	err := mq.Connect()
	if err != nil {
		logger.Info("MQ 生产者和消费者创建失败")
		return
	}
	mq.Subscribe("mq_request", "cn", func(topic, msgid string, body []byte) {
		logger.Info("top = ", topic)
		logger.Info("msgid = ", msgid)
		logger.Info("body = ", string(body))
	})
	msg := "hello mq"
	mq.Send("mq_request", []byte(msg), "cn")
}
