# 本地开发环境配置

:::tip 相关信息
以下内容将带领您逐步从零开始搭建和启动项目。如果您只想部署该项目，请转至[一键部署](/guide/deploy/docker)。
:::

## 后端环境

1. 前往 https://golang.google.cn/dl/ 下载对应系统的Go语言安装包。`推荐版本：1.20 +`
2. 命令行输入`go version`若控制台输出版本信息则安装成功
3. 使用命令 `go env -w GOPROXY=https://goproxy.cn,direct` 设置go代理


## 前端环境
1. 前往 https://nodejs.org/zh-cn/ 下载对应系统的Go语言安装包。`推荐版本：18.16 +`
2. 命令行输入`node -v`若控制台输出版本信息则安装成功

## 其他环境
:::tip 确保你的环境满足以下要求：
* Git: 推荐最新版本
* Mysql: >= 8.0，推荐使用8.0，其他版本未进行测试
* Redis: 推荐最新版本
* Nginx: 推荐最新版本
* FFmpeg: 推荐最新版本
:::
