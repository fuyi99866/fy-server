package websocket

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/service/robot_service"
	"net"
	"sync"
	"time"
)

type SubMsgInfo struct {
	RobotSN   string
	MessageID string
	Topic     string
	inputChan chan interface{}
	wsConn    *websocket.Conn
}

type WsClient struct {
	ws         *websocket.Conn        //当前websocket连接
	subMap     map[string]*SubMsgInfo //订阅表
	cancel     context.CancelFunc     //调用就断开websocket连接
	ctx        context.Context        //上下文
	Lock       *sync.RWMutex          //ws读写锁
	wsMsgChan  chan interface{}       //机器人消息通道，写入ws(因为ws写入不支持并发，所以使用select多路复用方式写入
	wsJsonChan chan interface{}       //机器人Json消息通道，写入ws
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (client *WsClient) IsMsgSubExist(RobotSN, Topic string, conn *websocket.Conn) bool {
	for _, v := range client.subMap {
		if v.RobotSN == RobotSN && v.Topic == Topic && v.wsConn == conn {
			return true
		}
	}
	return false
}

//创建机器人消息订阅
func (client *WsClient) AddMsgQueue(msgID, RobotSN, topic string, input chan interface{}, conn *websocket.Conn) error {
	client.subMap[msgID] = &SubMsgInfo{
		RobotSN:   RobotSN,
		MessageID: msgID,
		Topic:     topic,
		inputChan: input,
		wsConn:    conn,
	}
	logger.Info("AddMsgQueue", msgID, RobotSN, "current list", len(client.subMap))
	return nil
}

//取消机器人消息订阅
func (client *WsClient) RemoveMsgQueue(msgID string) {
	logger.Debug("RemoveMsgQueue", msgID)
	m := client.subMap[msgID]
	if m != nil {
		//TODO 取消对机器人消息的订阅
		robot_service.S.UnSubNotify(m.RobotSN,m.MessageID)
		delete(client.subMap, msgID)
	} else {
		logger.Error("RemoveMsgQueue error, has no", msgID)
	}
}

//断开连接
func (client *WsClient) Close() {
	close(client.wsMsgChan)
	close(client.wsJsonChan)
	for _, v := range client.subMap {
		//取消机器人消息的订阅
		robot_service.S.UnSubNotify(v.RobotSN,v.MessageID)
		delete(client.subMap, v.MessageID)
	}
	client.subMap = make(map[string]*SubMsgInfo, 0)
	client.ws.WriteMessage(websocket.CloseMessage, []byte{})
	client.ws.Close()
	logger.Debug("close")
}

//读取各种消息，然后通过websocket发送出去
func (client *WsClient) DispatchMessage() {
	defer func() {
		//捕获抛出的panic
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()
	for {
		select {
		case <-client.ctx.Done():
			logger.Info("quit DispatchMessage")
			return
		case dataStr, ok := <-client.wsMsgChan:
			if !ok {
				logger.Error("DispatchMessage error ,channel is not ok!", dataStr, ok)
				continue
			}
			err := client.ws.WriteMessage(websocket.TextMessage, dataStr.([]byte))
			if err != nil {
				logger.Error("websocket wsMsgChan error! maybe websocket has been closed, remove chan", string(dataStr.([]byte)), client.ws)
				client.cancel()
				return
			}
		case dataJson, ok := <-client.wsJsonChan:
			if !ok {
				logger.Error("DispatchMessage error ,channel is not ok!", dataJson, ok)
				continue
			}
			err := client.ws.WriteJSON(dataJson)
			if err != nil {
				logger.Error("WriteMessage  PingMessage err", err)
				client.cancel()
				return
			}
		}
	}
}

//处理web请求
func (client *WsClient) HandleWebRequest() {
	defer func() {
		//捕获抛出的panic
		if err := recover(); err != nil {
			logger.Warn(err)
		}
	}()
	for {
		select {
		case <-client.ctx.Done():
			logger.Info("quit HandleWebRequest")
			return
		default:
		}
		client.ws.SetReadDeadline(time.Now().Add(time.Minute * 3)) //if can't get read message after 3 min, websocket will close
		mt, data, err := client.ws.ReadMessage()
		if err != nil {
			client.cancel()
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					logger.Error("ReadMessage timeout remote: %v\n", client.ws.RemoteAddr())
				}
			} else {
				logger.Error("ReadWebsocketMessage err", client.ws.RemoteAddr(), err)
			}
			return
		}
		if string(data) == "ping" || mt == websocket.PingMessage {
			data = []byte("pong")
			client.wsMsgChan <- data
		} else if mt == websocket.PongMessage || string(data) == "pong" {
			logger.Debug("get websocket.PongMessage ")
		} else {
			logger.Debug("get websocket data", string(data), mt)
			var sub SubscribeRequest
			if err = json.Unmarshal(data, &sub); err != nil {
				client.wsJsonChan <- app.Response{
					Code: e.INVALID_PARAMS,
					Msg:  e.GetMsg(e.INVALID_PARAMS),
					Data: "",
				}
				continue
			}
			if sub.UnSubscribe { //取消订阅
				for k, v := range client.subMap {
					if v.RobotSN == sub.RobotSN && v.Topic == sub.Topic {
						client.RemoveMsgQueue(k)
						logger.Info("websocket UnSubscribe!, remove chan", v.RobotSN)
					}
				}
			} else { //订阅
				if client.IsMsgSubExist(sub.RobotSN, sub.Topic, client.ws) {
					logger.Debug("websocket SubNotify duplicate!", sub.RobotSN, sub.Topic)
					client.wsJsonChan <- app.Response{
						Code: e.ERROR_EXIST,
						Msg:  e.GetMsg(e.ERROR_EXIST),
						Data: sub.RobotSN,
					}
					continue
				}
				msgID, input, err := robot_service.S.SubNotify(sub.RobotSN)
				logger.Info("msgID, input, err : ", msgID, input, err)
				if err != nil {
					logger.Error("websocketSubscribeNotify error! ")
					client.wsJsonChan <- app.Response{
						Code: e.ERROR_NOT_EXIST,
						Msg:  e.GetMsg(e.ERROR_NOT_EXIST),
						Data: sub.RobotSN,
					}
					continue
				}else {
					client.AddMsgQueue(msgID,sub.RobotSN,sub.Topic,input,client.ws)
					go client.readRobotMsg(input,msgID,sub.RobotSN)
				}
			}
		}
	}
}

//读取机器人消息
func (client *WsClient) readRobotMsg(c chan interface{}, mid, sn string) {
	defer func() {
		err := recover()
		logger.Error("websocket read robot message error: ", err)
	}()
	for {
		select {
		case <-client.ctx.Done():
			logger.Info("quit readRobotMsg")
			return
		case data, ok := <-c:
			if !ok {
				logger.Error(mid, " ", data)
				return
			}
			if data == nil {
				logger.Error(mid, " ", data)
				return
			}
			logger.Info("readRobotMsg: ",data)
			client.wsMsgChan <- data
		}
	}
}
