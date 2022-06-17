参考：https://sjqzhang.github.io/go-fastdfs/ca.html



## 一、现象描述

在对fastdfs工程进行验证的时候，提示如下：

![111.png](https://tva1.sinaimg.cn/large/007Xg1efgy1h35sdgvb85j30ex0710ub.jpg)



而真实需要返回https格式的url





## 二、解决办法

### 2.1、更新版本

go-fastdfs文件服务器最新的版本才支持https，现在用的版本太旧了，还不支持，替换成新版本





### 2.2、修改配置文件

```bash
$ more cfg.json 
... ...
        "是否开启https": "默认不开启，如需启开启，请在conf目录中增加证书文件 server.crt 私钥 文件 server.key",
        "enable_https":true,
        
        "管理ip列表": "用于管理集的ip白名单,",
        "admin_ips": ["127.0.0.1","0.0.0.0"],
```



### 2.3、生成证书

进入到conf目录下



1.key的生成

```
openssl genrsa -des3 -out server.key 2048
openssl rsa -in server.key -out server.key
```



2. 生成CA的crt

```
openssl req -new -x509 -key server.key -out ca.crt -days 3650
```



3. csr的生成方法

```
openssl req -new -key server.key -out server.csr
```



4. crt生成方法

```
openssl x509 -req -days 3650 -in server.csr -CA ca.crt -CAkey server.key -CAcreateserial -out server.crt
```



### 2.4、重启fastdfs进程

```bash
 # nohup ./fileserver server &
```





## 2.5、修改nginx配置

将nginx关于fastdfs反向代理部分由http改成https

```bash
$ vi prerelease.ubtrobot.com.http 
        location ^~ /group1/ {
            proxy_set_header Host $host;
            proxy_redirect off;
            proxy_store off;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_connect_timeout 300s;
            proxy_send_timeout 300s;
            proxy_read_timeout 300s;
            proxy_pass https://dfs;   # 改成https
        }


```





## 

## 三、验证

![企业微信截图_20220613171757.png](https://tva1.sinaimg.cn/large/007Xg1efgy1h36qkmkcmbj30k90ae401.jpg)