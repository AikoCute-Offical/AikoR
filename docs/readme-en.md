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


# Description of AikoR
AikoR Supports Various Panels (V2board, ProxyPanel, sspanel, Pmpanel...)

An Xray-based back-end framework, supporting V2ay, Trojan, Shadowsocks protocols, extremely easily extensible and supporting multi-panel connection。

If you like this project, you can click the star + view in the upper right corner to track the progress of this project.

## Disclaimer

This project is for my personal learning, development and maintenance only, I do not guarantee the availability and I am not responsible for any consequences resulting from using this software.

## Featured
* Open source `This version depends on the happy mood`
* Supports multiple protocols V2ray, Trojan, Shadowsocks.
* Supports new features like Vless and XTLS.
* Supports single connection to multiple boards and nodes without rebooting.
* Online IP support is limited
* Support node port level, user level rate limit.
* Simple and clear configuration.
* Modify the configuration to automatically restart the instance.
* Easy to compile and upgrade, can quickly update core version, support new Xray-core features.
* Support UDP and many other functions

## Featured

| Featured                                       | v2ray | trojan | shadowsocks |
| -------------------------------------------    | ----- | ------ | ----------- |
| Get button info                                | √     | √      | √           |
| Get user information                           | √     | √      | √           |
| User traffic statistics                        | √     | √      | √           |
| Report server information                      | √     | √      | √           |
| Automatic registration of TLS certificates     | √     | √      | √           |
| auto-renew tls certificate                     | √     | √      | √           |
| Number of people online                        | √     | √      | √           |
| Online User Restrictions                       | √     | √      | √           |
| Audit rules                                    | √     | √      | √           |
| Node port speed limit                          | √     | √      | √           |
| User speed limit                               | √     | √      | √           |
| Custom DNS                                     | √     | √      | √           |
## User interface support

| Panel                                                  | v2ray | trojan | shadowsocks                                 |
| ------------------------------------------------------ | ----- | ------ | ------------------------------------------- |
| [sspanel-uim](https://github.com/Anankke/SSPanel-Uim)  | √     | √      | √ (Single-port multi-user and V2ray-Plugin) |
| [v2board](https://github.com/v2board/v2board)          | √     | √      | √                                           |
| [PMPanel](https://github.com/ByteInternetHK/PMPanel)   | √     | √      | √                                           |
| [ProxyPanel](https://github.com/ProxyPanel/ProxyPanel) | √     | √      | √                                           |

## Software installation - release
```
wget --no-check-certificate -O AikoR.sh https://raw.githubusercontent.com/AikoCute-Offical/AikoR-Install/master/AikoR.sh && bash AikoR.sh
```
### One main installation - docker
```
docker pull aikocute/aikor:latest && docker run --restart=always --name aikor -d -v ${PATCH_TO_CONFIG}/aiko.json:/etc/AikoR/aiko.json --network=host aikocute/aikor:latest
```
### Configuration file and detailed instructions
Comming Soon
## Telgram

Comming Soon

## Stargazers over time

[![Stargazers over time](https://starchart.cc/AikoCute-Offical/AikoR.svg)](https://starchart.cc/AikoCute-Offical/AikoR)
