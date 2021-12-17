package robot_service

import (
	"context"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"go_server/pkg/util"
	"sync"
	"time"
)

type MQ struct {
	url    string
	cli    mqtt.Client
	opts   *mqOptions
	mutex  sync.Mutex         //mqtt 锁
	ctx    context.Context    //上下文
	cancel context.CancelFunc //退出
}

func MakeMqttUrl(scheme, addr string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, addr, port)
}

func NewMq(url string, opts ...Option) *MQ {
	defaultOpts := defaultOptions()
	for _, apply := range opts {
		apply(&defaultOpts)
	}
	defaultOpts.url = url
	m := &MQ{
		opts: &defaultOpts,
	}
	conf := mqtt.NewClientOptions()
	if defaultOpts.tls {
		conf.SetTLSConfig(DefaultTlsConfig())
	}
	conf.SetClientID(defaultOpts.clientId)
	conf.AddBroker(url)
	conf.SetUsername(defaultOpts.username)
	conf.SetPassword(defaultOpts.password)
	conf.SetAutoReconnect(false) //自动重连
	conf.SetCleanSession(true)   //不清除订阅信息
	conf.SetConnectionLostHandler(m.mqttConnectionLostHandler)
	conf.SetOnConnectHandler(m.mqttOnConnectionHandler)
	conf.SetReconnectingHandler(m.mqttOnReconnectingHandler)
	c := mqtt.NewClient(conf)
	m.cli = c
	m.url = url
	m.ctx, m.cancel = context.WithCancel(context.Background())
	c.Connect()
	return m
}

//连接
func (m *MQ) Connect() error {
	logrus.Debug("Mqtt connect : ", m.opts, m.url)
	if m.cli.IsConnected() {
		return nil
	}
	token := m.cli.Connect()
	token.Wait()
	err := token.Error()

	if err != nil {
		return err
	}

	return nil
}

//重连
func (m *MQ) MustConnect(c context.Context) {
	util.RetryCancelWithContext(c, func() error {
		err := m.Connect()
		if err != nil {
			logrus.Info("Connect to mqtt  failed! ", err)
		} else {
			logrus.Info("Connect to mqtt  succeed!")
		}
		return err
	}, -1, time.Duration(m.opts.retry))
}

//断开
func (m *MQ) DisConnect() error {
	if m.cancel != nil {
		m.cancel()
	}
	if !m.cli.IsConnected() {
		return errors.New("mqtt connection has not been constructed")
	}
	m.cli.Disconnect(1000) //1000ms后断开
	return nil
}

//是否断开
func (m *MQ) IsConnected() bool {
	return m.cli.IsConnected()
}

//TODO 处理断开重连
func (m *MQ) mqttConnectionLostHandler(client mqtt.Client, err error) {
	logrus.Warnln("MQTT:Disconnected", err, client.IsConnected())
	if m.opts.autoReconnect {
		go m.MustConnect(m.ctx)
	}

	if m.opts.OnConnectionLost != nil {
		m.opts.OnConnectionLost(err)
	}
}

func (m *MQ) mqttOnConnectionHandler(mqtt.Client) {
	logrus.Info("Mqtt is connected")
	if m.opts.OnConnect != nil {
		m.opts.OnConnect()
	}
}

func (m *MQ) mqttOnReconnectingHandler(mqtt.Client, *mqtt.ClientOptions) {
	logrus.Info("Mqtt is reconnecting")
	if m.opts.OnReconnecting != nil {
		m.opts.OnReconnecting()
	}
}

//发送消息
func (m *MQ) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if topic == "" {
		return errors.New("publish topic can't be empty")
	}

	token := m.cli.Publish(topic, qos, retained, payload)
	go func(mqtt.Token) {
		_ = token.Wait()
		if token.Error() != nil {
			logrus.Println("Publish: ", token.Error())
		}
	}(token)
	return nil
}

//订阅
func (m *MQ) Subscribe(topic string, qos byte) error {
	if topic == "" {
		return errors.New("publish topic can't be empty")
	}
	token := m.cli.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		if m.opts.OnMessageComing != nil {
			m.opts.OnMessageComing(message.Topic(), message.Payload())
		} else {
			logrus.Infoln("There are no handler for handle", message.Topic())
		}
	})

	go func(mqtt.Token) {
		_ = token.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if token.Error() != nil {
			logrus.Println("UnSubscribe failed: ", token.Error())
		}
	}(token)
	return nil
}

//取消订阅
func (m *MQ) UnSubscribe(topic string) error {
	if topic == "" {
		return errors.New("publish topic can't be empty")
	}
	token := m.cli.Unsubscribe(topic)
	go func(mqtt.Token) {
		_ = token.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
		if token.Error() != nil {
			logrus.Println("UnSubscribe failed: ", token.Error())
		}
	}(token)
	return nil
}
