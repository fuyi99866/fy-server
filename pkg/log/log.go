package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"go_server/pkg/file"
	"go_server/pkg/setting"
	"log"
	"os"
	"time"
)

var logFilePath string
var logOutputs *LogOutputs
var logFile *os.File

func Init() {
	InitLog(setting.AppSetting.LogSavePath,
		setting.AppSetting.LogSaveName,
		setting.AppSetting.LogFileExt)
}

func InitLog(path, name, ext string) {
	LogFileName = name
	logFileExt = ext
	logPath = path

	logOutputs = NewLogOutputs()
	err := file.IsNotExistMkDir(path)
	if err != nil {
		logrus.Infoln(err)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	logOutputs.Appendlogrus(os.Stdout)
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		path+"/"+LogFileName+"%Y%m%d%H."+logFileExt,   //日志命名
		rotatelogs.WithMaxAge(time.Duration(24*3)*time.Hour),    //日志保留3天
		rotatelogs.WithRotationTime(time.Duration(1)*time.Hour), //每小时分包
	)
	logOutputs.Appendlogrus(writer)

	log.SetOutput(logOutputs)

	logrus.SetOutput(logOutputs)
	logrus.SetFormatter(&nested.Formatter{
		NoFieldsColors:  true, //warn、info、debug等颜色不同
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05", //time.RFC3339,
	})

	if setting.AppSetting.LogLever == "trace" {
		logrus.SetLevel(logrus.TraceLevel)
		gin.DefaultWriter = logOutputs
	} else if setting.AppSetting.LogLever == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
		gin.DefaultWriter = logOutputs
	} else {
		logrus.SetReportCaller(false)
		logrus.SetLevel(logrus.InfoLevel)
	}
}
