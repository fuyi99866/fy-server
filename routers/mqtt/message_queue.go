package mqtt

/**
消息队列
*/

//连接项
type ConnectOptions struct {
	Username       string
	Password       string
	MessageHandler MessageHandler
}

//消息处理器
type MessageHandler interface {
	//处理消息
	HandleMessage(string, []byte) error
	//处理断开连接
	HandleDisconnected(error)
}

type MessageQueue interface {
	Connect(ConnectOptions, string, string, bool) error
	Disconnect() error
	Publish(string, byte, bool, []byte) error
	Subscribe(string) error

	IsConnected() bool
}
