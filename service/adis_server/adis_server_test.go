package adis_server

import (
	"flag"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"testing"
)

func TestAll(t *testing.T) {
	logger.InitLog(setting.AppSetting.LogLever, "./logs/log.txt") //初始化日志库 ,使用zap库
	//config := flag.String("c", "D:/fy/git/go_server/fy-server/conf/app.ini", "config file path")
	config := flag.String("c", "D:/fy-server/go_server/conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置

	//初始化队列，创建生产者和消费者
	/*	logger.Info("setting.RocketMq.Addr: ", setting.RocketMq.Addr)
		logger.Info("setting.RocketMq.GroupName: ", setting.RocketMq.GroupName)*/
	/*	mq := NewMq(WithAddr(setting.RocketMq.Addr),
			WithGroupName(setting.RocketMq.GroupName),
			WithRetry(2),
		)
		err := mq.Connect()
		if err != nil {
			logger.Info("MQ 生产者和消费者创建失败")
			return
		}
		logger.Info("mq connect success")

		//TODO 订阅消息和发送消息还需要改一下
		mq.Subscribe("mq_request", "cn", func(topic, msgid string, body []byte) {
			logger.Info("topic = ", topic)
			logger.Info("msgid = ", msgid)
			logger.Info("body = ", string(body))
		})

		msg := "hello mq"
		mq.Send("mq_request", []byte(msg), "cn")

		time.AfterFunc(1000, func() {
			logger.Info("!!!!!!!!!!!!!")
		})*/

	/*	body, err := util.HttpPost("C:/Users/ubt/Desktop/45.jpg", "http://10.10.17.15:8087/group1/upload?scene=default&path=log&output=json")
		if err != nil {
			logger.Error("上传失败")
			return
		}

		logger.Info("body = ", string(body))*/

	//测试断点续传
	url, err := util.UploadByBreakPoint("C:/Users/ubt/Desktop/package-V1.2.9-2.zip", "http://10.10.17.15:8087/group1/big/upload/")
	if err != nil {
		return
	}
	logger.Info("文件下载路径 = ", url)

	body, err := util.HttpGetBySecond(url)
	if err != nil {
		logger.Info("err = ", err)
		return
	}

	logger.Info("body = ", string(body))

	//util.DownloadFileProgress("http://10.10.17.15:8087/group1/default/20220325/11/49/2/640c1c2ce613a201e92101a68b650bcb.zip","C:/Users/ubt/Desktop/1.zip")
}
