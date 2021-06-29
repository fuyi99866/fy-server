package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

//删除日志文件
func RemoveLogFile(logName string) {
	_, err := os.Stat(logName) //返回logName的文件信息
	if os.IsNotExist(err) {
		return
	}
	err_r := os.Remove(logName)
	if err_r != nil {
		logrus.Info("Fail to remove log file:", err_r)
	}else {
		logrus.Infoln("Success to remove log file:", logName)
	}
}
