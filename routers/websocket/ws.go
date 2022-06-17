package websocket

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

func Quit() {
	logrus.Info("socket quit")
	wsManager.ctx.Done()
}

// @Summary  Websocket接口, 支持订阅机器人状态、任务状态
// @Tags  websocket
// @Accept json
// @Produce  json
// @Param token query string true "TOKEN"
// @Param  body  body   SubscribeRequest  false "body"
// @Success 200 {object} MsgPackage
// @Failure 400 {object} app.Response
// @Router /channel [GET]
// @Security ApiKeyAuth
func NotifySocket(c *gin.Context) {
	defer func() {
		//捕获抛出的panic
		if err := recover(); err != nil {
			logrus.Info(err)
		}
	}()

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK,"create websocket failed")
		return
	}
	client:=new(WsClient)
	client.subMap = make(map[string]*SubMsgInfo,0)
	client.wsMsgChan = make(chan interface{},10)
	client.wsJsonChan = make(chan interface{},10)
	client.ctx,client.cancel = context.WithCancel(wsManager.ctx)
	client.ws = ws

	//向前端发送消息
	go client.DispatchMessage()
	//处理请求消息
	go client.HandleWebRequest()

	select {
	case <-client.ctx.Done()://等待退出
	}
	client.Close()
}
