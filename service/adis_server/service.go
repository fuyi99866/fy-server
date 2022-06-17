package adis_server

import (
	"go_server/pkg/logrus"
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
		logrus.Info("mq connect failed: ", err)
		return err
	}
	logrus.Info("mq connect success")

	//订阅MQ消息
	a.Subscribe("mq_request", func(topic, msgid string, body []byte) {
		logrus.Info("rocketmq message coming:", topic, " ",msgid, " ",string(body))
	})

	time.AfterFunc(1000, func() {
		msg:="hello mq"
		a.Send("mq_request",[]byte(msg))
	})

	return nil
}

func TestStart()  {
	//初始化队列，创建生产者和消费者
	logrus.Info("setting.RocketMq.Addr: ", setting.RocketMq.Addr)
	logrus.Info("setting.RocketMq.GroupName: ", setting.RocketMq.GroupName)
	mq := NewMq(WithAddr(setting.RocketMq.Addr),
		WithGroupName(setting.RocketMq.GroupName),
		WithRetry(2),
	)
	err := mq.Connect()
	if err != nil {
		logrus.Info("MQ 生产者和消费者创建失败")
		return
	}
	logrus.Info("mq connect success")

	//TODO 订阅消息和发送消息还需要改一下
	mq.Subscribe("mq_request", "cn", func(topic, msgid string, body []byte) {
		logrus.Info("topic = ", topic)
		logrus.Info("msgid = ", msgid)
		logrus.Info("body = ", string(body))
	})

	msg := "hello mq"
	mq.Send("mq_request", []byte(msg), "cn")
}

//结束服务
func (a *adisService) Stop() error {
	a.mq.DisConnect()
	return nil
}

//发送消息
func (a *adisService) Send(topic string, msg []byte) {
	logrus.Info("adisService Send :", topic, " ",string(msg))
	a.mq.Send(topic, msg, "cn")
}

//订阅消息
func (a *adisService) Subscribe(topic string, f ReceiveCallBack) {
	logrus.Info("adisService Subscribe topic :", topic)
	a.mq.Subscribe(topic, "cn", f)
}
