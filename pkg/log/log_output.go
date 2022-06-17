package log

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type LogOutputs struct {
	logruss  []io.Writer
	logFlags int
}

func NewLogOutputs() *LogOutputs {
	var outputs LogOutputs
	outputs.logruss = make([]io.Writer, 0)
	return &outputs
}

func (outputs *LogOutputs) Appendlogrus(writer io.Writer) error {
	if outputs.logruss == nil {
		return errors.New("LogOutputs is nil")
	}
	if writer == os.Stdout {
		outputs.logFlags = log.Flags()
	}

	outputs.logruss = append(outputs.logruss, writer)
	return nil
}

func (outputs *LogOutputs) Removelogrus(writer io.Writer) error {
	if outputs.logruss == nil {
		return errors.New("outputs.logruss is nil")
	}

	var index = -1
	for i, v := range outputs.logruss {
		if v == writer {
			index = i
		}
	}

	if index != -1 {
		outputs.logruss = append(outputs.logruss[:index], outputs.logruss[index+1:]...)
		return nil
	}

	return errors.New("writer not in range")
}

func (outputs *LogOutputs) Write(p []byte) (n int, err error) {
	if outputs.logruss == nil {
		return 0, errors.New("outputs.logruss is nil")
	}
	for _, v := range outputs.logruss {
		/*		var msg = string(p)
				var newMsg = msg*/
		n, err = v.Write(p)
	}
	return n, err
}

var LogFileName string
var logFileExt string
var logPath string

func GetLogFileName(date time.Time) string {
	return path.Join(logPath+"/", LogFileName+date.Format("20060102")+"."+logFileExt)
}

func RemoveLogFile(logName string) {
	_, err := os.Stat(logName)
	if os.IsNotExist(err) {
		return
	}
	err = os.Remove(logName)
	if err != nil {
		logrus.Infoln("Fail to remove log file:", err)
	} else {
		logrus.Infoln("Success to remove log file:", logName)
	}
}
