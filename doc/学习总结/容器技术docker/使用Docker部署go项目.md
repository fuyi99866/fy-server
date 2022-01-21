#### 使用docker部署go项目
* 1、下载安装docker
```cassandraql
sudo yum install -y yum-utils device-mapper-persistent-data lvm2 
sudo yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo 
sudo yum install docker
sudo systemctl enable docker
sudo systemctl start docker
```
* 2、编写Dockerfile文件
```cassandraql
FROM ubuntu:18.04

WORKDIR /go_server
COPY ./go_server /go_server/go_server
COPY ./dist /go_server/dist
COPY ./conf/app.ini /go_server/conf/app.ini
RUN chmod +x /go_server/go_server

EXPOSE 8082

CMD ./go_server -c ./conf/app.ini
```
* 3、生成docker镜像
```cassandraql
docker build -t go-server-docker .
```
* 4、运行docker镜像
```cassandraql
docker run -p 8082:8082 go-server-docker
```
* 5、docker常用指令
##### 拉取镜像
```cassandraql
docker pull 
```
##### 查看镜像
```cassandraql
docker images
```
##### 查询正在运行的容器
```cassandraql
docker ps
docker ps -a # 为查看所有的容器，包括已经停止的
```
##### 删除容器
```cassandraql
docker rm <容器名 or ID>
```
##### 删除所有容器
```cassandraql
docker rm $(docker ps -a -q)
```
##### 停止、启动、杀死指定容器
```cassandraql
docker start <容器名 or ID> # 启动容器
docker stop <容器名 or ID> # 启动容器
docker kill <容器名 or ID> # 杀死容器
```
##### 查看所有镜像
```cassandraql
docker images
```
##### 拉取镜像
```cassandraql
docker pull <镜像名:tag>
```
##### 删除镜像
```cassandraql
docker rmi <镜像ID>
```
##### 端口暴露
```cassandraql
# 一共有三种形式进行端口映射
docker -p ip:hostPort:containerPort # 映射指定地址的主机端口到容器端口
# 例如：docker -p 127.0.0.1:3306:3306 映射本机3306端口到容器的3306端口
docker -p ip::containerPort # 映射指定地址的任意可用端口到容器端口
# 例如：docker -p 127.0.0.1::3306 映射本机的随机可用端口到容器3306端口
docer -p hostPort:containerPort # 映射本机的指定端口到容器的指定端口
# 例如：docker -p 3306:3306 # 映射本机的3306端口到容器的3306端口
```

#### 使用docker-compose部署go项目
* 1、编写docker-compose文件
```text
version: '3'
services:
  go_server:
    build: .
    ports:
    - "8082:8082"
  redis:
    image: redis
  emqx:
    image: emqx/emqx
    ports:
    - 18083:18083
    - 1883:1883
    - 8883:8883
```
* 2、启动项目
```text
docker-compose up
或者
docker-compose up -d #后台运行
```
* 3、常用指令
##### 查看服务
```text
docker-compose ps
```
##### 启动、停止服务
```text
docker-compose start [name]
docker-compose stop [name]
```
##### 删除服务
```text
docker-compose rm [name]
```
##### 查看具体服务的日志
```text
docker-compose logs -f [name]
```
##### 可以进入容器内部
```text
docker-compose exec [name] shell
```

##### 参考资料
[使用docker部署一个go应用](https://www.cnblogs.com/ricklz/p/12860434.html)

[Gin实践 连载九 将Golang应用部署到Docker](https://segmentfault.com/a/1190000013960558)

[如何使用Docker部署一个Go Web应用程序](http://dockone.io/article/1269)