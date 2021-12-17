package robot_service

import (
	"github.com/sirupsen/logrus"
	"go_server/pkg/util"
	"sync"
)

type RobotOnlineChangedHandler func(cid, sn string, online bool) error
type RobotStatusChangedHandler func(cid, sn string) error

//定义单个机器人的数据
type ROBOT struct {
	SN            string
	Company       string
	chanMutex     sync.Mutex //同步锁
	notifyChanMap map[string]chan interface{}
	connect       bool
	online        bool
	mq            *MQ
}

func NewROBOT(sn, cid string, mq *MQ) *ROBOT {
	r := &ROBOT{
		SN:            sn,
		Company:       cid,
		chanMutex:     sync.Mutex{},
		notifyChanMap: make(map[string]chan interface{}),
		online:        false,
		mq:            mq,
	}
	return r
}

//初始化单个机器人的数据
func (robot *ROBOT) initTopic() error {
	logrus.Info("ROBOT init topic :", robot.SN)
	if err := robot.mq.Subscribe(MakeTopic(robot.Company, robot.SN, ROBOT_REQUEST), 0); err != nil {
		logrus.Info("Subscribe failed :", err)
		robot.connect = false
		return err
	}

	if err := robot.mq.Subscribe(MakeTopic(robot.Company, robot.SN, ROBOT_RESPONSE), 0); err != nil {
		logrus.Info("Subscribe failed :", err)
		robot.connect = false
		return err
	}

	if err := robot.mq.Subscribe(MakeTopic(robot.Company, robot.SN, ROBOT_NOTIFY), 0); err != nil {
		logrus.Info("Subscribe failed :", err)
		robot.connect = false
		return err
	}

	if err := robot.mq.Subscribe(MakeTopic(robot.Company, robot.SN, ROBOT_CONNECT), 0); err != nil {
		logrus.Info("Subscribe failed :", err)
		robot.connect = false
		return err
	}

	if err := robot.mq.Subscribe(MakeTopic(robot.Company, robot.SN, ROBOT_DISCONNECT), 0); err != nil {
		logrus.Info("Subscribe failed :", err)
		robot.connect = false
		return err
	}

	robot.connect = true
	logrus.Info("ROBOT init topic success:", robot.SN)
	return nil
}

//退出时取消订阅
func (robot *ROBOT) Quit() {
	robot.mq.UnSubscribe(MakeTopic(robot.Company, robot.SN, ROBOT_REQUEST))
	robot.mq.UnSubscribe(MakeTopic(robot.Company, robot.SN, ROBOT_RESPONSE))
	robot.mq.UnSubscribe(MakeTopic(robot.Company, robot.SN, ROBOT_NOTIFY))
	robot.mq.UnSubscribe(MakeTopic(robot.Company, robot.SN, ROBOT_CONNECT))
	robot.mq.UnSubscribe(MakeTopic(robot.Company, robot.SN, ROBOT_DISCONNECT))
	msgStr := CreateOnlineInfo(robot.SN, false) //通知前端离线
	robot.PubNotify(msgStr)
	for mid, _ := range robot.notifyChanMap {
		robot.UnSubNotify(mid)
	}
}

//更新企业号
func (robot *ROBOT) Update(cid string) {
	robot.Quit()
	robot.Company = cid
	robot.initTopic()
}

//订阅通知
func (robot *ROBOT) SubNotify() (string, chan interface{}, error) {
	logrus.Infoln("Robot SubNotify ", robot.SN)
	robot.chanMutex.Lock()
	defer robot.chanMutex.Unlock()
	messageId := util.UUIDShort()
	logrus.Infoln("Robot SubNotify messageId", robot.SN, messageId)
	c := make(chan interface{}, 10)
	robot.notifyChanMap[messageId] = c //存储messageId
	return messageId, c, nil
}

//取消订阅
func (robot *ROBOT) UnSubNotify(messageId string) {
	logrus.Infoln("Robot UnSubNotify", robot.SN, messageId)
	robot.chanMutex.Lock()
	defer robot.chanMutex.Unlock()
	close(robot.notifyChanMap[messageId]) //close message chan
	delete(robot.notifyChanMap, messageId)
}

//推送消息
func (robot *ROBOT) PubNotify(data []byte) {
	if data == nil || len(data) == 0 {
		return
	}
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorln("PubNotify send robot message error:", err)
		}
	}()
	logrus.Debugln("send robot msg:", robot.SN, string(data), len(robot.notifyChanMap))
	robot.chanMutex.Lock()
	defer robot.chanMutex.Unlock()
	for id, c := range robot.notifyChanMap {
		logrus.Debugln("send robot msg to web", id)
		go util.SafeSend(c, data)
	}
}
