package robot_service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_server/models"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"strings"
	"sync"
	"time"
)

/**
建立MQTT通道，与机器人连接
1、增加机器人的时候、开启连接
2、启动的时候读取机器人列表，建立连接
*/

//所有机器人管理
type robotService struct {
	mq                  *MQ
	robots              sync.Map
	msgList             sync.Map //request&response
	connecting          bool
	cancel              context.CancelFunc //重连取消
	RobotOnlineCallback RobotOnlineChangedHandler
	RobotStatusCallback RobotStatusChangedHandler
}

func NewService() *robotService {
	s := &robotService{
		connecting: false,
	}
	return s
}

//消息队列通用格式
type MsgPackage struct {
	Result    chan string
	MessageId string
	Topic     string
	HoldTitle string
	Cancel    context.Context
}

var S *robotService = NewService()

//开启机器人服务
func (s *robotService) Start() error {
	s.mq = NewMq(MakeMqttUrl(setting.Mqtt.Scheme, setting.Mqtt.Addr, setting.Mqtt.Port),
		WithUserName(setting.Mqtt.UserName),
		WithPassword(setting.Mqtt.Password),
		WithClientId(fmt.Sprintf("ROBOT_%s", util.GetRandomString(10))),
		WithTls(setting.Mqtt.Tls),
		WithReconnectingHandler(s.onReconnecting),
		WithConnectionLostHandler(s.onDisconnected),
		WithConnectedHandler(s.onConnected),
		WithMessageHandler(s.onMessageComing))
	s.connectMqtt()
	return nil
}

func (s *robotService) onMessageComing(topic string, data []byte) {
	isOnlineTopic := strings.HasPrefix(topic, "$SYS/brokers/")
	if !isOnlineTopic {
		cid, sn, scheme, _ := ParseTopic(topic)
		if scheme == ROBOT_REQUEST {
			s.onRequest(cid, sn, data)
		}
		if scheme == ROBOT_RESPONSE {
			s.onResponse(cid, sn, data)
		}
		if scheme == ROBOT_NOTIFY {
			s.onNotify(cid, sn, data)
		}
	} else {
		sn, scheme, _ := parseOnlineTopic(topic)
		online := true
		if scheme == ROBOT_CONNECT {
			online = true
		} else if scheme == ROBOT_DISCONNECT {
			online = false
		}

		logrus.Info("sn , online :", sn, online)
		//TODO 保存机器人的在线状态，并向前端推送
	}
}

func (s *robotService) onReconnecting() {
	logrus.Warnln("robotService:onReconnecting", s.connecting)
}

func (s *robotService) onDisconnected(err error) {
	logrus.Warnln("robotService:onDisconnected", err, s.connecting)
}

func (s *robotService) onConnected() {
	logrus.Warnln("robotService:onConnected  ", "connected ", s.mq.IsConnected(), "connecting ", s.connecting)
}

//结束机器人服务
func (s *robotService) Stop() error {
	logrus.Warnln("Stop robotService")
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.mq.DisConnect()
	return nil
}

//连接机器人
func (s *robotService) ConnectRobot(cid, sn string) error {
	return nil
}

//断开机器人
func (s *robotService) DisconnectRobot(cid, sn string) error {
	return nil
}

//检验机器人是否存在
func (s *robotService) existRobot(sn string) bool {
	rvalue, exist := s.robots.Load(sn)
	return (rvalue != nil) && exist
}

//获取机器人信息
func (s *robotService) getRobot(sn string) *ROBOT {
	rvalue, exist := s.robots.Load(sn)
	if rvalue == nil || !exist {
		return nil
	}
	robot, _ := rvalue.(*ROBOT)
	return robot
}

//收到请求
func (s *robotService) onRequest(cid, sn string, data []byte) error {
	return nil
}

//收到回复
func (s *robotService) onResponse(cid, sn string, data []byte) error {
	return nil
}

//收到通知
func (s *robotService) onNotify(cid, sn string, data []byte) error {
	return nil
}

//连接到mqtt
func (s *robotService) connectMqtt() {
	if s.connecting {
		return
	}
	s.connecting = true
	s.robots.Range(func(key, value interface{}) bool {
		robot, _ := value.(*ROBOT)
		if robot == nil {
			return true
		}
		s.DisconnectRobot(robot.Company, robot.SN)
		return true
	})

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	//先断开，后重连

	go util.RetryCancelWithContext(ctx, func() error {
		err := s.mq.Connect()
		if err != nil {
			s.connecting = true
		} else {
			s.connecting = false
			logrus.Infoln("Connect to  MQTT" + " succeed")
			robots, _ := models.GetAllRobot()
			for _, robot := range robots {
				s.ConnectRobot(robot.Company, robot.SN)
			}
		}
		return err
	}, -1, time.Second*3) //3秒重连一次，无限重连
}
