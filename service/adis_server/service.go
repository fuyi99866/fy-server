package adis_server

import (
	"github.com/sirupsen/logrus"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"time"
)

//定义服务接口
type AdisService interface {
	Start() error
	Stop() error
	Send(topic string, msg []byte)
	Subscribe(topic string, f ReceiveCallBack)
}

//定义服务结构体
type adisService struct {
	mq *MQ
	//MsgList *[]interface{}
}

func NewService() *adisService {
	a := &adisService{}
	return a
}

var A = NewService()

//开启服务
func (a *adisService) Start() error {
	a.mq = NewMq(
		WithAddr(setting.RocketMq.Addr),
		WithGroupName(setting.RocketMq.GroupName),
		WithRetry(2),
	)
	err := a.mq.Connect()
	if err != nil {
		logger.Info("mq connect failed: ", err)
		return err
	}
	logger.Info("mq connect success")

	//订阅MQ消息
	a.Subscribe("mq_request", func(topic, msgid string, body []byte) {
		logrus.Debugln("rocketmq message coming:", topic, msgid, string(body))
	})

	time.AfterFunc(1000, func() {
		msg:="hello mq"
		a.Send("mq_request",[]byte(msg))
	})

	return nil
}

//结束服务
func (a *adisService) Stop() error {
	a.mq.DisConnect()
	return nil
}

//发送消息
func (a *adisService) Send(topic string, msg []byte) {
	logger.Info("adisService Send :", topic, string(msg))
	a.mq.Send(topic, msg, "cn")
}

//订阅消息
func (a *adisService) Subscribe(topic string, f ReceiveCallBack) {
	logger.Info("adisService Subscribe topic :", topic)
	a.mq.Subscribe(topic, "cn", f)
}
