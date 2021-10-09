###### 在ubuntu登录数据库，并进行操作-u 
####### Password = c72d988fe3f10ed394a7888b9645bb9f
#######Host = 172.31.0.58:3306
```text
mysql -h host地址 -u 用户名 -p
show databases;
use 
select * from adis_robot;
```

###### 使用命令让机器人连上公司网络
####### ubuntu账号:cruiser@10.5.5.1
####### 密码:ubt123
```text
sudo /home/cruiser/wifi/wifi_ctl.sh sta '"UBT-Users"' '"ubtubtubt"'
```
