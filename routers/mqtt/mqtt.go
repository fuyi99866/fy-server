package mqtt

import (
	"crypto/tls"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type MessageQueueMQTT struct {
	client         mqtt.Client
	MessageHandler MessageHandler
	mutex          sync.Mutex
}

func newTlsConfig() *tls.Config {
	return &tls.Config{
		Rand:                        nil,
		Time:                        nil,
		Certificates:                nil,
		NameToCertificate:           nil,
		GetCertificate:              nil,
		GetClientCertificate:        nil,
		GetConfigForClient:          nil,
		VerifyPeerCertificate:       nil,
		VerifyConnection:            nil,
		RootCAs:                     nil,//broker ca
		NextProtos:                  nil,
		ServerName:                  "",
		ClientAuth:                  tls.NoClientCert,
		ClientCAs:                   nil,
		InsecureSkipVerify:          false, //skip verify the server cert
		CipherSuites:                nil,
		PreferServerCipherSuites:    false,
		SessionTicketsDisabled:      false,
		SessionTicketKey:            [32]byte{},
		ClientSessionCache:          nil,
		MinVersion:                  0,
		MaxVersion:                  0,
		CurvePreferences:            nil,
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}
}

//判断MQTT的连接状态
func (mq *MessageQueueMQTT) IsConnected() bool {
	if mq.client == nil {
		return false
	}
	return mq.client.IsConnected()
}

//处理MQTT断开事件
func (mq *MessageQueueMQTT) ConnectionLostHandler(client mqtt.Client, err error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	if mq.MessageHandler != nil {
		mq.MessageHandler.HandleDisconnected(err)
	}
}

//MQTT 连接
func (mq *MessageQueueMQTT) Connect(options ConnectOptions, url string, clientID string, tls bool) error {
	var mqttClientOptions = mqtt.NewClientOptions()
	mqttClientOptions.SetClientID(clientID)
	mqttClientOptions.AddBroker(url)
	mqttClientOptions.SetUsername(options.Username)
	mqttClientOptions.SetPassword(options.Password)
	mqttClientOptions.SetConnectTimeout(time.Second*2)//设置连接超时时间为2秒
	mqttClientOptions.SetAutoReconnect(false)//设置为不自动重连
	//处理断开事件
	mqttClientOptions.SetConnectionLostHandler(mq.ConnectionLostHandler)
	if tls{
		tlsConfig:=newTlsConfig()
		mqttClientOptions.SetTLSConfig(tlsConfig)
	}
	mq.client = mqtt.NewClient(mqttClientOptions)
	mq.MessageHandler = options.MessageHandler
	if mq.client == nil{
		return errors.New("failed to connect to")
	}
	token:=mq.client.Connect()//连接服务器
	token.Wait()
	err:=token.Error()
	if err!=nil{
		return err
	}
	return nil
}

//发布消息
func (mq *MessageQueueMQTT) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if topic == "" || !mq.IsConnected() {
		return errors.New("Invalid topic or connection")
	}
	token := mq.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

//订阅
func (mq *MessageQueueMQTT) Subscribe(topic string, qos byte) error {
	if topic == "" || !mq.IsConnected() {
		return errors.New("Invalid topic or connection ")
	}
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	token := mq.client.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		mq.mutex.Lock()
		defer mq.mutex.Unlock()
		if mq.MessageHandler != nil {
			mq.MessageHandler.HandleMessage(message.Topic(), message.Payload())
		} else {
			logrus.Info("MessageHandler is nil")
		}
	})

	token.Wait()
	return token.Error()
}

//取消订阅
func (mq *MessageQueueMQTT) UnSubscribe(topic string) error {
	if topic == "" || !mq.IsConnected() {
		return errors.New("Invalid topic or connection ")
	}
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	token := mq.client.Unsubscribe(topic)
	token.Wait()
	return token.Error()
}

//断开 MQTT
func (mq *MessageQueueMQTT) Disconnect() error {
	if !mq.IsConnected() {
		return errors.New("connection is not exist")
	}
	mq.client.Disconnect(100)
	return nil
}
