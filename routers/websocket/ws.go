package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
)

//前端订阅websocket消息
type SubscribeRequest struct {
	RobotSN     string `json:"robotSn"`
	Topic       string `json:"topic"`       //default：Notify
	UnSubscribe bool   `json:"unsubscribe"` //true：取消订阅
}

//http请求转为websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//返回值
type MsgPackage struct {
	RobotSN string
	Type    string
	Data    string
}

type WsManager struct {
	ctx context.Context
}

var wsManager *WsManager = new(WsManager)

func init() {
	wsManager.ctx = context.Background()
}

func Quit()  {
	wsManager.ctx.Done()
}
