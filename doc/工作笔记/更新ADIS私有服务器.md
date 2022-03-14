###### 更新ADSI私有服务器
####### 密码: adis123
```text
ssh adis.ubtrobot.com
cd ADIS/ADISS-Server/
git pull
go build
sudo systemctl stop adis     
sudo cp -r adis-server /usr/bin/adis/adis-server
sudo systemctl restart adis
```

###### 建勇的数据库
```text
冼建勇 9-22 16:15:21
[app]
PageSize = 10
JwtSecret = 1234567890
#linux
LogSavePath = /var/log/adis
#windows
#LogSavePath = logfile/adis/
LogSaveName = adis
LogFileExt = log
#LogLever debug or info
LogLever = debug
TimeFormat = 20060102

[server]
#RunMode debug or release
RunMode = debug
HttpPort = 9090
ReadTimeout = 60
WriteTimeout = 60
HTTPS = false
BasePath = /v1/adis

[database]
Type = mysql
TablePrefix = adis_

[sqlite]
Version = sqlite3
Path = ./adis.db

[mysql]
User = adis
Password = c72d988fe3f10ed394a7888b9645bb9f
Host = 172.31.0.58:3306
Name = adis
TablePrefix = adis_

[swagger]
Host = prerelease.ubtrobot.com

[JWT]
#hour
ExpiresTime = 24

[MQTT]
#"tcp", "ssl", or "ws"
Scheme = ssl
Addr = prerelease.ubtrobot.com
Port = 16666
RobotUser = adis
RobotPassword = c72d988fe3f10ed394a7888b9645bb9f
Tls = true

冼建勇 9-23 10:54:51
[文件：2021-09.zip]
```

##### 生产环境深圳节点MQTT配置文件
```
Scheme = ssl
Addr = adis.ubtrobot.com
Port = 18888
RobotUser = adis
RobotPassword = c72d988fe3f10ed394a7888b9645bb9f
Tls = true
```
##### 生产环境新加坡节点MQTT配置文件
```
Scheme = ssl
Addr = adis-sg.ubtrobot.com
Port = 16666
RobotUser = adis
RobotPassword = c72d988fe3f10ed394a7888b9645bb9f
Tls = true
```

##### 生产环境北美节点MQTT配置文件
```text
Scheme = ssl
Addr = adis-na.ubtrobot.com
Port = 16666
RobotUser = adis
RobotPassword = c72d988fe3f10ed394a7888b9645bb9f
Tls = true
```
