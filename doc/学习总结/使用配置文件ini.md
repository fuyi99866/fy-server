### 使用配置文件ini
#### 编写项目配置包
##### 拉取go-ini/ini
```text
go get -u github.com/go-ini/ini
```
##### 新建conf文件夹，在里面创建app.ini文件写入下面内容
```text
#debug or release
RUN_MODE = debug

[app]
PAGE_SIZE = 10
JWT_SECRET = 23347$040412

[server] 
HTTP_PORT = 8081
READ_TIMEOUT = 60
WRITE_TIMEOUT = 60

[database]
TYPE = mysql
USER = root
PASSWORD = root
#127.0.0.1:3306
HOST = localhost:8081
NAME fy
TABLE_PREFIX = fy_
```

##### 建立调用配置的setting模块
