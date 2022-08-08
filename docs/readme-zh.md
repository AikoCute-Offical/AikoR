<p align="center"><img src="https://avatars.githubusercontent.com/u/91626055?v=4" width="128" /></p>

<div align="center">

# AikoR
AikoR Projects

[![](https://img.shields.io/badge/Telegram-group-green?style=flat-square)](https://t.me/AikoXrayR)
[![](https://img.shields.io/badge/Telegram-channel-blue?style=flat-square)](https://t.me/AikoCute_Support)
[![](https://img.shields.io/github/downloads/AikoCute-Offical/AikoR/total.svg?style=flat-square)](https://github.com/AikoCute-Offical/AikoR/releases)
[![](https://img.shields.io/github/v/release/AikoCute-Offical/AikoR?style=flat-square)](https://github.com/AikoCute-Offical/AikoR/releases)
[![docker](https://img.shields.io/docker/v/aikocute/aikor?label=Docker%20image&sort=semver)](https://hub.docker.com/r/aikocute/aikor)
[![Go-Report](https://goreportcard.com/badge/github.com/AikoCute-Offical/AikoR?style=flat-square)](https://goreportcard.com/report/github.com/AikoCute-Offical/AikoR)
</div>

# AikoR的描述
AikoR 支持各种面板（V2board、ProxyPanel、sspanel、Pmpanel...）

基于Xray的后端框架，支持V2ay、Trojan、Shadowsocks协议，极易扩展，支持多面板连接。

如果你喜欢这个项目，可以点击右上角的星号+视图来跟踪这个项目的进度。

## 免责声明

本项目仅供我个人学习、开发和维护，我不保证可用性，也不对使用本软件产生的任何后果负责。

## 精选
* 开源`此版本取决于心情愉快`
* 支持多种协议V2ray、Trojan、Shadowsocks。
* 支持 Vless 和 XTLS 等新功能。
* 支持单连接多板多节点，无需重启。
* 在线 IP 支持有限
* 支持节点端口级别、用户级别限速。
* 简单明了的配置。
* 修改配置自动重启实例。
* 易于编译和升级，可以快速更新核心版本，支持新的Xray-core特性。
* 支持UDP等多项功能

## 精选

|精选 | v2ray |木马 |影袜|
| ------------------------------------------------------- | ----- | ------ | ------------ |
|获取按钮信息 | √ | √ | √ |
|获取用户信息 | √ | √ | √ |
|用户流量统计 | √ | √ | √ |
|报告服务器信息 | √ | √ | √ |
|自动注册 TLS 证书 | √ | √ | √ |
|自动更新 tls 证书 | √ | √ | √ |
|在线人数 | √ | √ | √ |
|在线用户限制 | √ | √ | √ |
|审计规则 | √ | √ | √ |
|节点端口限速 | √ | √ | √ |
|用户限速 | √ | √ | √ |
|自定义 DNS | √ | √ | √ |
##用户界面支持

|面板 | v2ray |木马 |影袜|
| -------------------------------------------------- ---- | ----- | ------ | ------------------------------------------------------- |
| [sspanel-uim](https://github.com/Anankke/SSPanel-Uim) | √ | √ | √（单端口多用户和V2ray-Plugin）|
| [v2board](https://github.com/v2board/v2board) | √ | √ | √ |
| [PMPanel](https://github.com/ByteInternetHK/PMPanel) | √ | √ | √ |
| [代理面板](https://github.com/ProxyPanel/ProxyPanel) | √ | √ | √ |

## 软件安装-发布
```
wget --no-check-certificate -O AikoR.sh https://raw.githubusercontent.com/AikoCute-Offical/AikoR-Install/master/AikoR.sh && bash AikoR.sh
```
### 一个主要安装 - docker
```
docker pull aikocute/aikor:latest && docker run --restart=always --name aikor -d -v ${PATCH_TO_CONFIG}/aiko.json:/etc/AikoR/aiko.json --network=host aikocute/aikor:latest
```
### 配置文件和详细说明
即将推出
##电报

即将推出

## 随时间推移的观星者