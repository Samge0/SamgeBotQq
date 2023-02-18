## QQ机器人简易助手
基于 [go-Pichubot](https://github.com/AkiinuBot/go-Pichubot) 二次开发；
数据交互基于 [gocp](https://docs.go-cqhttp.org/) 的正向WebSocket；
使用docker快捷运行，方便在qq中使用openai体验[chatGPT](https://chat.openai.com/chat)~

【温馨提示】非官方产品有被封号风险，请使用小号作为Bot尝试。

### 基本步骤

 - 1、注册openai并获取其token；
 - 2、运行gocp作为正向websocket服务端；
 - 3、运行本项目代码进行消息处理；

### 运行gocp作为正向websocket服务端
在ubuntu中使用docker运行gocp

```shell
sudo groupadd -g 1024 qqbot
sudo useradd -M -g qqbot -s /sbin/nologin -u 1024 qq123456789
```

运行docker镜像，其中6700为本项目在config.yml中配置的正向websocket端口，5700是http的端口
```docker
docker run -it \
--name qqbot \
-v /home/samge/docker_data/qqbot:/data \
-p 5700:5700 \
-p 6700:6700 \
-e TZ="Asia/Shanghai" \
-e UID=1024 \
-e GID=1024 \
ghcr.io/mrs4s/go-cqhttp:1.0.0-rc4
```
第一次运行会在映射目录下生成config.yml配置文件（具体配置介绍参考 [go-cqhttp的配置](https://docs.go-cqhttp.org/guide/config.html#%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF) ），在配置文件中根据需要，启用http或者正向websocket的配置项，
填写qq、qq密码、token、http的端口、正向websocket的地址跟端口，比如本项目仅使用正向websocket，则只需要配置正向websocket的信息即可。
配置完毕后保存，重新运行docker镜像（记得-p将配置的端口映射到物理机）。
编写后端程序，链接正向websocket，监听消息，处理消息，到此已完成qq消息的建议机器人，如果需要其他深度定制功能，需自己根据需求修改代码或使用其他更成熟的框架（参考后面链接）。

### 运行本项目代码进行消息处理

 - `-v xxx:/app/tmp/qqBotCache` =》这里的xxx填写自己的映射路径，存放日志，例如下面的`/home/samge/docker_data/samge_qq_bot`
 - `OPEN_AI_TOKEN=`填写openai的token值，
 - `GROUP_IDS=`填写群的白名单，多个用,分隔（qq群号）
 - `FRIEND_IDS=`填写好友的白名单，多个用,分隔（好友qq号）
 - `MINE_NICKNAME=`填写当前机器人bot的昵称
 - `MASTER_QQ=`填写管理者qq，多个用,分隔
 - `SOCKET_HOST=`填写qq通讯的正向websocket地址
 - `SOCKET_TOKEN=`填写与websocket通讯的token（gocp中配置的token）

```shell
docker run -d \
--name samge_qq_bot \
-v /home/samge/docker_data/samge_qq_bot:/app/tmp/qqBotCache \
--pull=always \
--restart always \
-e LANG=C.UTF-8 \
-e TZ="Asia/Shanghai" \
-e OPEN_AI_TOKEN= \
-e GROUP_IDS= \
-e FRIEND_IDS= \
-e MASTER_QQ= \
-e SOCKET_HOST= \
-e SOCKET_TOKEN= \
samge/samge_qq_bot:v1
```

### 如果需要调试

 - 配置：`cmd/qqBot/botConfig/botConfig.go`
 - 安装依赖：`go mod tidy`
 - 运行：`go run cmd/qqBot/main.go`

### 类似的聊天机器人消息处理框架有：

 - [ZeroBot（GO）-可同时多个bot](https://github.com/wdvxdr1123/ZeroBot)
 - [nonebot2（Python）](https://github.com/nonebot/nonebot2)
 - [nonebot机器人商店（Python）](https://v2.nonebot.dev/store)
 - [OneBot（多语言）](https://onebot.dev/ecosystem.html#%E5%BA%94%E7%94%A8%E6%A1%88%E4%BE%8B)

### 其他

 - [基于Nonebot2和go-cqhttp的机器人搭建](https://yzyyz.top/archives/nb2b1.html)
 - [go-cqhttp的配置](https://docs.go-cqhttp.org/guide/config.html#%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF)
 - [国内聊天机器人调研](https://www.jianshu.com/p/a1ee997b1330)
 - [nb跟gocp的介绍](https://yzyyz.top/archives/nb2b1.html)


### 免责声明
该程序仅供技术交流，使用者所有行为与本项目作者无关
