##### 本机服务器调试ADIS接口
* 1、连接机器人的WIFI，如"977A"
* 2、打开本地的静态uv_tool调试页面
* 3、在调试页面“主机”上填写：10.5.5.1，并登陆和验证
* 4、设置服务器，将IP和端口换成本机的IP和端口
* 5、获取机器人信息，查询机器人的SN
* 6、将SN添加到ADIS

##### 重置机器人的同步进度
```text
//更新数据同步进度
sqlite3 ~/impdata/uv_database.db
update uv_sync_progress_temp_tb set sync_id=0;
//重启节点
//查看进度是否归零
select *from  uv_sync_progress_temp_tb ;
.quit
rosnode kill uv_business_unit_node;
```