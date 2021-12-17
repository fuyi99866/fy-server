package robot_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_server/pkg/util"
	"go_server/service/robot_service/format"
	"strings"
	"time"
)

const (
	ROBOT_REQUEST  = "Request"
	ROBOT_RESPONSE = "Response"
	ROBOT_NOTIFY   = "Notify"

	ROBOT_CONNECT    = "connected"
	ROBOT_DISCONNECT = "disconnected"
)

//创建MQTT的topic
func MakeTopic(cid, sn, topic string) string {
	logrus.Debugln("MakeTopic", fmt.Sprintf("%s/%s/%s", cid, sn, topic))
	return fmt.Sprintf("%s/%s/%s", cid, sn, topic)
}

//创建MQTT检测机器人上下线的topic
func MakeOnlineTopic(sn, topic string) string {
	logrus.Debug("MakeOnlineTopic: ", fmt.Sprintf("$SYS/brokers/+/clients/%s/%s", sn, topic))
	return fmt.Sprintf("$SYS/brokers/+/clients/%s/%s", sn, topic)
}

//从topic获取上下线信息
func parseOnlineTopic(topic string) (sn, scheme string, err error) {
	isOnlineTopic := strings.HasPrefix(topic, "$SYS/brokers/")
	values := strings.Split(topic, "/")
	logrus.Debugln("ParseOnlineTopic : ", topic, isOnlineTopic, values)
	if !isOnlineTopic || len(values) < 6 {
		return "", "", errors.New("error topic")
	}
	return values[4], values[5], nil
}

//从topic获取sn、企业号，消息类型等信息
func ParseTopic(topic string) (cid, sn, scheme string, err error) {
	values := strings.Split(topic, "/")
	if len(values) != 3 {
		return "", "", "", errors.New("error topic")
	}
	return values[0], values[1], values[2], nil
}

func CreateOnlineInfo(sn string, online bool) []byte {
	resp := format.WebSocketResponseNotifyOnline{
		Title: format.RobotOnlineStateNotifyTitle,
		Content: struct {
			SessionID string `json:"sessionid"`
			Timestamp int64  `json:"timestamp"`
			Data      struct {
				Sn     string `json:"sn"`
				Online bool   `json:"online"`
			} `json:"data"`
		}{
			SessionID: util.UUIDShort(),
			Timestamp: time.Now().Unix(),
			Data: struct {
				Sn     string `json:"sn"`
				Online bool   `json:"online"`
			}{
				Sn:     sn,
				Online: online,
			},
		},
	}
	msgStr, _ := json.Marshal(resp) //结构体转为json
	return msgStr
}
