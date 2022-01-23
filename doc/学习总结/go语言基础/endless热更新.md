#### 原理：
主进程接受到 SIGTERM 信号量，fork 了新的子进程，已经发出的请求在旧进程处理，新的请求在新的进程处理，直到所有请求都由新进程处理时，关闭旧进程
#### 代码实现
```cassandraql
package main

import (
    "fmt"
    "log"
    "syscall"

    "github.com/fvbock/endless"

    "gin-blog/routers"
    "gin-blog/pkg/setting"
)

func main() {
	r := gin.New()
    endless.DefaultReadTimeOut = setting.ReadTimeout
    endless.DefaultWriteTimeOut = setting.WriteTimeout
    endless.DefaultMaxHeaderBytes = 1 << 20
    endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

    server := endless.NewServer(endPoint, r )
    server.BeforeBegin = func(add string) {
        log.Printf("Actual pid is %d", syscall.Getpid())
    }

    err := server.ListenAndServe()
    if err != nil {
        log.Printf("Server err: %v", err)
    }
}
```
#### linux环境验证指令
```cassandraql
# 第一次构建项目
go build main.go
# 运行项目，这时就可以做内容修改了
./endless &
# 请求项目，60s后返回
curl "http://127.0.0.1:5003/sleep?duration=60s" &
# 再次构建项目，这里是新内容
go build main.go
# 重启，17171为pid
kill -1 17171
# 新API请求
curl "http://127.0.0.1:5003/sleep?duration=1s" 
```