package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

//定义传输的数据结构体
type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
}

//用户结构体
type User struct {
	ws        *websocket.Conn //当前websocket连接
	wsMsgChan chan []byte     //消息通道，写入ws
	data      *Data           //发送的数据信息，未封装
}

//数据处理器
type Hub struct {
	userList   map[*User]bool //用户列表,保存所有的用户
	broadcast  chan []byte    //消息通道
	register   chan *User     //注册chan，用户注册时添加到chan中
	unregister chan *User     //注销chan，用户退出时添加到chan中，再从map中删除
}

//定义四个通道
var h = Hub{
	userList:   make(map[*User]bool), //用户组映射
	register:   make(chan *User),     //用户加入通道
	broadcast:  make(chan []byte),    //消息通道
	unregister: make(chan *User),     //用户退出通道
}

var user_list = []string{}


var upGrader = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary  Websocket接口, 支持订阅机器人状态、任务状态、任务倒计时信息
// @Accept json
// @Produce  json
// @Tags  websocket
// @Accept json
// @Produce  json
// @Router /ws [GET]
func NotifySocket(c *gin.Context) {
	go h.Run()
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("获取连接失败:", err)
		return
	}
	//建立websocket连接
	user := &User{wsMsgChan: make(chan []byte, 256), ws: ws, data: &Data{}}

	//用户表中没有该用户，则用户加入注册通道
	h.register <- user

	defer func() {
		h.unregister <- user
	}()

	//得到连接后开始读写数据
	go user.SendMsg()
	user.ReadMsg()

}

//处理中心处理获取到的信息
func (h *Hub) Run() {
	for {
		select {
		//从注册chan中取数据
		case user := <-h.register:
			h.userList[user] = true
			user.data.Ip = user.ws.RemoteAddr().String()
			user.data.Type = "handshake"
			user.data.UserList = user_list
			data_b, _ := json.Marshal(user.data)
			logrus.Println("data_b == ", string(data_b))
			user.wsMsgChan <- data_b
		//从注销列表中取数据，判断用户列表中是否存在这个用户，存在就删掉
		case user := <-h.unregister:
			if _, ok := h.userList[user]; ok {
				delete(h.userList, user)
				close(user.wsMsgChan)
			}
		//从广播chan中取消息，然后遍历给每个用户，发送到用户的wsMsgChan中
		case data := <-h.broadcast:
			for c := range h.userList {
				select {
				case c.wsMsgChan <- data:
				default:
					delete(h.userList, c)
					close(c.wsMsgChan)
				}
			}
		}
	}
}

//发送消息
func (user *User) SendMsg() {
	for message := range user.wsMsgChan {
		logrus.Println("发送信息： ", string(message))
		user.ws.WriteMessage(websocket.TextMessage, message)
	}
	user.ws.Close()
}

//读取消息
func (user *User) ReadMsg() {
	for {
		_, message, err := user.ws.ReadMessage()
		logrus.Debugln("get websocket message", string(message))
		if err != nil {
			logrus.Errorln("用户退出： ", user.ws.RemoteAddr().String())
			h.unregister <- user
			break
		}

		//message消息json反序列化为结构体
		json.Unmarshal(message, &user.data)
		switch user.data.Type {
		case "login":
			user.data.User = user.data.Content
			user.data.From = user.data.User
			user_list = append(user_list, user.data.User)
			user.data.UserList = user_list
			data_b, _ := json.Marshal(user.data) //将数据编码成json字符串
			h.broadcast <- data_b
			logrus.Println("user.data.User == ",user.data.User)
		case "user":
			user.data.Type = "user"
			data_b, _ := json.Marshal(user.data)
			h.broadcast <- data_b
			logrus.Println("user:data_b == ",string(data_b))
		case "logout":
			user.data.Type = "logout"
			user_list = del(user_list, user.data.User)
			data_b, _ := json.Marshal(user.data)
			h.broadcast <- data_b
			h.register <- user
		default:
			fmt.Print("============default+++++++++++++++")
		}

	}
}

//删除用户
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println("=======", n_slice)
	return n_slice
}
