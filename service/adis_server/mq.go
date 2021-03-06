package adis_server

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"go_server/pkg/logrus"
	"golang.org/x/net/context"
)

//MQ 使用rocketmq消息中间件，与adis进行通讯
type MQ struct {
	producer rocketmq.Producer
	consumer rocketmq.PushConsumer
	opts     *mqOptions
}

//初始化日志
func init() {
	rlog.SetLogLevel("error")
}

//NewMq 创建一个MQ客户端
func NewMq(opts ...Option) *MQ {
	var m = new(MQ)
	defaultOpts := defaultOptions()
	for _, apply := range opts {
		apply(&defaultOpts)
	}
	m.opts = &defaultOpts
	logrus.Debug("adis NewMq ", m.opts.addr, defaultOpts.addr)
	return m
}

//创建生产者和消费者
func (m *MQ) Connect() error {
	err := m.createProducer()
	if err != nil {
		return err
	}

	err = m.createConsumer()
	if err != nil {
		return err
	}
	return nil
}

//shutdown生产者和消费者
func (m *MQ) DisConnect() error {
	if m.producer != nil {
		m.producer.Shutdown()
	}
	if m.consumer != nil {
		m.consumer.Shutdown()
	}
	return nil
}

//发送topic的异步消息
func (m *MQ) Send(topic string, msg []byte, tag string) {
	//发送异步消息
	err := m.producer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			logrus.Error("send rcmp message : ", err)
		} else {
			logrus.Info("send rcmp message success. result= ", result.String())
		}
	}, primitive.NewMessage(topic, msg).WithTag(tag))
	if err != nil {
		logrus.Error("send adis message : ", err)
	}
}

func (m *MQ) Subscribe(topic, tag string, f ReceiveCallBack) {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	} //只订阅带Tag的消息
	if tag == "" {
		selector = consumer.MessageSelector{}
	}
	err := m.consumer.Subscribe(topic, selector, func(ctx context.Context, ext ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
		for _, v := range ext {
			t := v.GetTags()
			logrus.Info("adis Subscribe callback tag: ", t)
			go f(v.Topic, v.MsgId, v.Body) //没懂?
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		logrus.Error("Subscribe adis message : ", err)
	}
}

//创建生产者
func (m *MQ) createProducer() error {
	addr, err := primitive.NewNamesrvAddr(m.opts.addr)
	if err != nil {
		logrus.Error("createProducer NewNamesrvAddr failed : ", err.Error())
		return err
	}
	p, err := rocketmq.NewProducer(
		producer.WithGroupName(m.opts.groupName),
		producer.WithNameServer(addr),
		producer.WithRetry(m.opts.retry),
	)
	if err != nil {
		logrus.Error("createProducer NewProducer failed : ", err.Error())
		return err
	}
	if err = p.Start(); err != nil {
		logrus.Error("createProducer Start failed : ", err.Error())
		return err
	}
	m.producer = p
	return nil
}

//创建消费者
func (m *MQ) createConsumer() error {
	//消息主动推送给消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithInstance(m.opts.consumerInstance), //必须设置，否则广播模式会重复消费
		consumer.WithGroupName(m.opts.groupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{m.opts.addr})),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), //选择消费时间（首次/当前/根据时间）
		consumer.WithConsumerModel(consumer.Clustering),                //消费模式(集群消费:消费完,同组的其他人不能再读取/广播消费：所有人都能读)
	)
	if err != nil {
		logrus.Error("createConsumer NewPushConsumer failed : ", err.Error())
		return err
	}
	err = c.Start()
	if err != nil {
		logrus.Error("createConsumer Start failed : ", err.Error())
		return err
	}
	m.consumer = c
	return nil
}
