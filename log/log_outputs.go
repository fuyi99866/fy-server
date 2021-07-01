package log

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	DebugTextFormatter = "\033[0;32m[debug]\033[0m" // green text color
	WarnTextFormatter  = "34m[warn]"                // blue text color
	BeginColorGroup    = "\033[0;"
	EndColorGroup      = "\033[0m"
)

const (
	MAX_PREFIX_LEN = 75
)

const (
	BLACK  = "30m"
	RED    = "31m"
	GREEN  = "32m"
	YELLOW = "33m"
	BLUE   = "34m"
	PURPLE = "35m"
	CYAN   = "36m"
	WHITE  = "37m"
)

type LogOutputs struct {
	loggers  []io.Writer
	logFlags int
}

func NewLogOutputs() *LogOutputs {
	var outputs LogOutputs
	outputs.loggers = make([]io.Writer, 0)
	return &outputs
}

func (outputs *LogOutputs) AppendLogger(writer io.Writer) error {
	if outputs.loggers == nil {
		return errors.New("AppendLogger error")
	}

	//如果是标准输出
	if writer == os.Stdout {
		outputs.logFlags = log.Flags()
	}

	outputs.loggers = append(outputs.loggers, writer)
	return nil
}

func (outputs *LogOutputs) RemoveLogger(writer io.Writer) error {
	if outputs.loggers == nil {
		return errors.New("RemoveLogger error")
	}
	var index = -1
	for i, v := range outputs.loggers {
		if v == writer {
			index = i
		}
	}

	if index != -1 {
		outputs.loggers = append(outputs.loggers[:index], outputs.loggers[index+1:]...)
		return nil
	}

	return errors.New("writer is not in range")
}

func (outputs *LogOutputs) Writer(p []byte) (n int, err error) {
	if outputs.loggers == nil {
		return 0, errors.New("can not write log")
	}

	for _, v := range outputs.loggers {
		var msg = string(p)
		var newMsg = msg
		if v == os.Stdout { //标准输出
			var logFlags = outputs.logFlags
			var msgs = strings.Split(msg, " ")
			var index = 0
			if (logFlags & log.Ldate) != 0 {
				newMsg = BeginColorGroup + GREEN + msgs[index] + EndColorGroup + " "
				index++
				newMsg += BeginColorGroup + PURPLE + msgs[index] + EndColorGroup + " "
				index++
			}

			if (logFlags & log.Lshortfile) != 0 {
				newMsg += BeginColorGroup + YELLOW + msgs[index] + EndColorGroup + " "
				msgIndex := strings.Index(msg, msgs[index]) + len(msgs[index]) + 1
				newMsg += msg[msgIndex:]
			} else {
				newMsg += EndColorGroup + msg
			}
		}
		n, err = v.Write([]byte(newMsg))
	}

	return n, err
}

var LogFileName string
var logFileExt string
var logPath string

//获取日志文件名
func GetLogFileName(date time.Time) string {
	return path.Join(logPath+"/", LogFileName+date.Format("20060102")+"."+logFileExt)
}

//删除日志文件
func RemoveLogFile(logName string) {
	_, err := os.Stat(logName) //返回logName的文件信息
	if os.IsNotExist(err) {
		return
	}
	err_r := os.Remove(logName)
	if err_r != nil {
		logrus.Info("Fail to remove log file:", err_r)
	} else {
		logrus.Infoln("Success to remove log file:", logName)
	}
}
