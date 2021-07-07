package setting

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"go_server/pkg/logger"
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

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

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

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		logger.Fatalf("加载配置文件失败", err)
	}
}

