package setting

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"time"
)

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type SWAGGER struct {
	Host string
}

var Swag = &SWAGGER{}

type SERVER struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	HTTPS        bool
	BasePath     string
}

var ServerSetting = &SERVER{}

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	PrefixUrl      string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	LogLever    string
	TimeFormat  string
}

var AppSetting = &App{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type MQTT struct {
	Scheme   string
	Addr     string
	Port     int
	UserName string
	Password string
	Tls      bool
}

var Mqtt = &MQTT{}

type ROCKETMQ struct {
	Addr string
	GroupName string
	Enable bool
}
var RocketMq = &ROCKETMQ{}

var cfg *ini.File
var RunMode string

func Init(config string) {
	var err error
	cfg, err = ini.Load(config)
	if err != nil {
		logrus.Fatal("初始化配置文件失败： ", err)
	}
	mapTo("app", AppSetting)
	mapTo("database", DatabaseSetting)
	mapTo("swagger", Swag)
	mapTo("redis", RedisSetting)
	mapTo("server", ServerSetting)
	mapTo("mqtt", Mqtt)
	mapTo("rocketmq",RocketMq)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		logrus.Fatalf("加载配置文件失败", err)
	}
}
