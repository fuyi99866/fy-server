package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

// logrus提供了New()函数来创建一个logrus的实例。
// 项目中，可以创建任意数量的logrus实例。
//var LOG = logrus.New()
var (
	MyLogger *myLogger
)

type myLogger struct {
	*logrus.Logger
	File *os.File
}

// logrus在记录Levels()返回的日志级别的消息时会触发HOOK，
// 按照Fire方法定义的内容修改logrus.Entry
type myHook struct {
	Field  string         //输出日志的代码文件名称
	Skip   int            //skip为遍历调用栈开始的索引位置
	levels []logrus.Level //日志等级
}

//实现 logrus.Hook 接口
//每次有日志消息写入时，会查询findCaller
func (hook *myHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

//实现 logrus.Hook 接口
func (hook *myHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//自定义hook
func NewMyHook(levels ...logrus.Level) logrus.Hook {
	hook := myHook{
		Field:  "field",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

//自定义logger
func NewLogger(level logrus.Level, format logrus.Formatter, hook logrus.Hook) *logrus.Logger {
	log := logrus.New()
	log.Level = level
	log.SetFormatter(format)
	log.Hooks.Add(hook)
	return log
}

func findCaller(skip int) string {
	file := ""
	line := 0
	var pc uintptr
	// 遍历调用栈的最大索引为第11层.
	for i := 0; i < 11; i++ {
		file, line, pc = getCaller(skip + i)
		//fmt.Println("findCaller", file, line)
		//文件名不能以logrus开头 ,过滤掉所有logrus包，即可得到生成代码信息
		if strings.HasSuffix(file, "logrus") {
			break
		}
	}

	fullFnName := runtime.FuncForPC(pc)
	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		//取得函数名
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d:%s()", file, line, fnName)
}

//当前goroutine调用栈中的文件名，行号，函数信息等，参数skip表示表示返回的栈帧的层次
func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	//fmt.Println("getCaller", pc, file, line, ok)
	if !ok {
		return "", 0, pc
	}
	n := 0
	//获取包名
	for i := len(file) - 1; i > 0; i-- {
		if string(file[i]) == "/" {
			n++
			if n >= 2 {
				//fmt.Println(n >= 2, file)
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}

//初始化配置
func Init() {
	var (
		file *os.File
		err  error
	)
	path := "log.txt"
	if file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
		logrus.Error("打开日志文件错误：", err)
	}
	MyLogger = &myLogger{
		File: file,
	}
	MyLogger.Logger = NewLogger(logrus.DebugLevel, &logrus.TextFormatter{FullTimestamp: true}, NewMyHook())
	MyLogger.Logger.Out = MyLogger.File
}
