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
sudo docker run --name nginx-test -p 8083:80 -d nginx
```
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
##### 部署
##### 在这里需要对 nginx.conf 进行配置，如果你不知道对应的配置文件是哪个，可执行 nginx -t 看一下
##### 配置hosts
由于需要用本机作为演示，因此先把映射配上去，打开 /etc/hosts，增加内容：
```cassandraql
127.0.0.1       api.blog.com
```
##### 配置 nginx.conf
打开 nginx 的配置文件 nginx.conf（我的是 /usr/local/etc/nginx/nginx.conf），我们做了如下事情：
增加 server 片段的内容，设置 server_name 为 api.blog.com 并且监听 8081 端口，将所有路径转发到 http://127.0.0.1:8000/ 下
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
        listen       8081;
        server_name  api.blog.com;

        location / {
            proxy_pass http://127.0.0.1:8000/;
        }
    }
}
```
##### 验证
启动服务
##### 重启 nginx
```cassandraql
nginx -t
```
##### 访问接口
```cassandraql
http://api.blog.com:8081/auth?username=admin&password=123456
```

#### 参考内容：
[docker配合Nginx部署go应用](https://www.cnblogs.com/ricklz/p/12996174.html)

[Golang Gin实践 连载十七 用 Nginx 部署 Go 应用](https://segmentfault.com/a/1190000016236253)

[vue+go-gin+nginx实现后台管理系统](https://blog.csdn.net/hahachenchen789/article/details/105847926/)

[Nginx反向代理+go后端服务](https://www.cnblogs.com/guichenglin/p/12760698.html)