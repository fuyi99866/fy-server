###### 在ubuntu登录数据库，并进行操作-u 
####### Password = c72d988fe3f10ed394a7888b9645bb9f
#######Host = 172.31.0.58:3306
```text
mysql -h host地址 -u 用户名 -p
show databases;
use adis;
select * from adis_robot;
```

###### 使用命令让机器人连上公司网络
####### ubuntu账号:cruiser@10.5.5.1
####### 密码:ubt123
```text
```

##### websocket调试
```text
wss://prerelease.ubtrobot.com/v1/adis/channel?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJpZCI6OCwicGFzc3dvcmQiOiIyNWQ1NWFkMjgzYWE0MDBhZjQ2NGM3NmQ3MTNjMDdhZCIsImNvbXBhbnlpZCI6IlRFU1QiLCJleHAiOjE2MzI5MTQzMzEsImlzcyI6Imh0dHBzOi8vYWRpcy1zZXJ2ZXIvIn0.ApAj4l0PjteEdYgKkwil-ZfHZ2h71iV9asKDGsjMOUA

ws://10.10.17.15:9090/v1/channel?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJpZCI6MiwicGFzc3dvcmQiOiJlMTBhZGMzOTQ5YmE1OWFiYmU1NmUwNTdmMjBmODgzZSIsImNvbXBhbnlpZCI6InVidGVjaCIsImV4cCI6MTYzMjg4NTgyOSwiaXNzIjoiaHR0cHM6Ly9hZGlzLXNlcnZlci8ifQ.6zPw6zL4-1jhzGHT8cSnxVEGYib1ck2hnV5IJ0VAtOU
```

```text
{
   	"robotsn": "CAI001UBT10000016",
   	"topic": "notify",
   	"unsubscribe": false
   }
```

##### 设置私有服务器地址
```text
http://10.10.17.15:9090/v1/profile
```


```text
wss://adis.ubtrobot.com/v1/adis/channel?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJpZCI6NDgsInBhc3N3b3JkIjoiMjVkNTVhZDI4M2FhNDAwYWY0NjRjNzZkNzEzYzA3YWQiLCJjb21wYW55aWQiOiJURVNUIiwiZXhwIjoxNjM3MDYzNzIzLCJpc3MiOiJodHRwczovL2FkaXMtc2VydmVyLyJ9.k__hP12DP9nf_0i70DpDJOpXaMCUK5PE455GM868idw
```

##### 查看adis预发布环境的日志
```text
cd /data/ubt/adis
tail -f -n 100 nohup.out
或者
cd /var/log/adis
tail -f -n 100 adis2022020911.log
```

##### 查看adis生产环境（深圳节点）的日志
```text
cd /data/ubt/adis/logs
tail -f -n 100 adis2022020911.log
```