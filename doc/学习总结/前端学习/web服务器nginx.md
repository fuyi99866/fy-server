#### Nginx部署go应用
##### Nginx是什么，用来做什么？
##### Nginx是一个 Web Server,可以用作反向代理、负载均衡、邮件代理、TCP/UDP、HTTP服务器等。
* 以较低的内存占用率处理10000多个并发连接
* 静态服务器（处理静态文件）
* 正向、反向代理
* 负载均衡
* 通过OpenSSL对TSL/SSL与SNI和OCSP支持
* FastCGI、SCGI、UWSGI的支持
* WebSockets、HTTP/1.1的支持
* Nginx+Lua
##### 安装
##### 使用docker安装
##### 拉取镜像
```cassandraql
docker pull nginx
```
##### 启动Nginx服务
```cassandraql
sudo docker run --name nginx -p 8083:80 -d nginx
```
##### 修改docker容器中的nginx的配置文件
* 进入容器
```cassandraql
docker exec -it nginx bash
```
* 进入配置文件目录
```cassandraql
cd /etc/nginx/
```
* 安装vim
```cassandraql
apt-get  update
apt-get install vim
```
* 打开配置文件
```cassandraql
vim conf.d/default.conf
```
* 修改配置文件
```cassandraql
    server {
        listen       80;
        server_name  localhost;

        location / {
            proxy_pass http://1.117.171.11:8081/swagger/index.html;
        }
    }
```
* 关闭nginx,再次启动
```cassandraql
docker stop nginx
docker start nginx
```
##### 或者使用挂载配置文件的方式装有docker宿主机上面的nginx.conf配置文件映射到启动的nginx容器里面，这需要你首先准备好nginx.con配置文件
```cassandraql
mkdir nginx
docker cp nginx-t:/etc/nginx/nginx.conf /home/ubuntu/nginx/
docker cp nginx-t:/etc/nginx/conf.d /home/ubuntu/nginx/
```
##### 修改配置文件conf.d/default.conf
```cassandraql
worker_processes  1;

events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

    server {
        listen       80;
        server_name  localhost;

        location / {
            proxy_pass http://1.117.171.11:8081/index.html;
        }
    }
}
```
##### 然后启动nginx
##### 命令：
```cassandraql
docker run --name nginx -p 80:80 -v /home/ubuntu/nginx/nginx.conf:/etc/nginx/nginx.conf -v /home/ubuntu/nginx/log:/var/log/nginx -v /home/ubuntu/nginx/conf.d/default.conf:/etc/nginx/conf.d/default.conf -d nginx
```
##### 解释下上面的命令：
--name  给你启动的容器起个名字，以后可以使用这个名字启动或者停止容器

-p 映射端口，将docker宿主机的80端口和容器的80端口进行绑定

-v 挂载文件用的，第一个-v 表示将你本地的nginx.conf覆盖你要起启动的容器的nginx.conf文件，第二个表示将日志文件进行挂载，就是把nginx服务器的日志写到你docker宿主机的/home/docker-nginx/log/下面

第三个-v 表示的和第一个-v意思一样的。

-d 表示启动的是哪个镜像

##### 常用指令
* nginx：启动Nginx
* nginx -s stop：立即停止Nginx服务
* nginx -s reload:重新加载配置文件
* nginx -s quit:平滑停止Nginx服务
* nginx -t:测试配置文件是否正确
* nginx -v：显示Nginx版本信息
* nginx -V：显示 Nginx 版本信息、编译器和配置参数的信息
##### 涉及配置
1、proxy_pass:配置反向代理路径。如果proxy_pass的url最后为/，则表示绝对路径。否则（不含变量下）表示相对路径，所有路径都会被代理过去
2、upstream:配置负载均衡，upstream默认是轮询的方式进行负载，另外还支持四种模式：
* weight：权重，指定轮询的概率，weight与访问概率成正比
* ip_hash:按照访问IP的hash结果值分配
* fair：按后端服务器响应时间进行分配，响应时间越短优先级越高
* url_hash：按照访问URL的hash结果分配


#### 参考内容：
[docker配合Nginx部署go应用](https://www.cnblogs.com/ricklz/p/12996174.html)

[Golang Gin实践 连载十七 用 Nginx 部署 Go 应用](https://segmentfault.com/a/1190000016236253)

[docker上启动nginx并修改nginx配置文件](https://www.cnblogs.com/hl15/p/13686515.html)

[nginx实现前后端分离](https://blog.csdn.net/mybook201314/article/details/88743861)