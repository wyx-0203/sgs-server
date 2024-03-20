# 基于Unity与Go语言的网络卡牌游戏

# 的设计与实现

2019024216 王宇曦            指导老师：何斌

## 简介

本项目使用Unity引擎复刻了经典卡牌桌游《三国杀》，包含单机模式与联网模式，实现了《三国杀》的核心游戏逻辑(包括多人回合制、卡牌系统、武将系统等)，以及用户系统、房间系统、联机对战等功能。游戏以三国时期为背景，以武将为角色，以卡牌为形式，合纵连横，经过一轮一轮的谋略和动作获得最终的胜利。

## 预览

<img src="images/preview1.png" style="zoom: 50%;" />

<img src="images/preview2.png" style="zoom:50%;" />

## 技术栈

#### 客户端:

* Unity
* C#

#### 服务端:

* Go

* Gin

* GORM

* WebSocket

#### 服务器相关:

* Docker

* Nginx

#### 数据库:

* MySQL

## 快速启动

### 环境要求

* Go (1.18+)

* Docker

* Docker Compose

### 安装说明

运行本应用将会启动3个Docker容器，分别是:

* `web-api`容器: 运行后端服务(本机编译Go项目，得到可执行文件，并传入容器)

* `nginx`容器: 代理443端口、部署静态文件(WebGL)、配置ssl

* `mysql`容器: 数据库


### 运行应用

```sh
make up
```

此命令会调用docker-compose构建并运行容器。

### 停止运行

```sh
make down
```

## 主要结构

```
sgs-server
├── controllers 控制器层，接收http请求
│   ├── match.go 游戏服务，如建立WebSocket连接/进入房间/开始游戏等，会调用match包
│   ├── personal.go 个人信息服务
│   └── user.go 用户服务
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── global
│   └── constants.go 配置文件
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── match 游戏服务
│   ├── hub.go 管理所有房间和在线玩家
│   ├── message.go 定义所有WebSocket消息结构体
│   ├── player.go 玩家结构体，用于保存信息和WebSocket连接
│   └── room.go 房间结构体
├── models 模型层，与数据库交互
│   ├── init.go
│   ├── personal.go 个人信息(昵称、形象)
│   └── user.go 用户
├── nginx
│   ├── cert SSL证书
│   ├── default.conf nginx配置文件
│   ├── default_local.conf nginx配置文件(本地调试版)
│   ├── Dockerfile
│   ├── index.html 网站主页(WebGL)
│   └── webgl Unity_WebGL静态文件
├── README.md
└── utils
    └── jwt.go JWT认证
```
