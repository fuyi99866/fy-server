package format
const (
	RobotOnlineStateNotifyTitle   = "notify_robot_online_state"   //通知前端机器人上下线状态
)

type WebSocketResponseNotifyOnline struct {
	Title   string `json:"title"`
	Content struct {
		SessionID string `json:"sessionid"`
		Timestamp int64  `json:"timestamp"`
		Data      struct {
			Sn     string `json:"sn"`
			Online bool   `json:"online"`
		} `json:"data"`
	} `json:"content"`
}