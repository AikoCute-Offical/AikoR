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

# Mô tả của AikoR
AikoR hỗ trợ nhiều bảng điều khiển khác nhau (V2board, ProxyPanel, sspanel, Pmpanel ...)

Một khung back-end dựa trên Xray, hỗ trợ các giao thức V2ay, Trojan, Shadowsocks, cực kỳ dễ dàng mở rộng và hỗ trợ kết nối nhiều bảng điều khiển。

Nếu bạn thích dự án này, bạn có thể nhấp vào dấu sao + xem ở góc trên bên phải để theo dõi tiến độ của dự án này.

## Tuyên bố từ chối trách nhiệm

Dự án này chỉ dành cho việc học tập, phát triển và bảo trì của cá nhân tôi, tôi không đảm bảo tính khả dụng và tôi không chịu trách nhiệm về bất kỳ hậu quả nào do sử dụng phần mềm này.

## Đặc sắc
* Mã nguồn mở `Phiên bản này phụ thuộc vào tâm trạng vui vẻ`
* Hỗ trợ nhiều giao thức V2ray, Trojan, Shadowsocks.
* Hỗ trợ các tính năng mới như Vless và XTLS.
* Hỗ trợ một kết nối đến nhiều bo mạch và nút mà không cần khởi động lại.
* Hỗ trợ IP trực tuyến bị hạn chế
* Hỗ trợ mức cổng nút, giới hạn tốc độ mức người dùng.
* Cấu hình đơn giản và rõ ràng.
* Sửa đổi cấu hình để tự động khởi động lại phiên bản.
* Dễ dàng biên dịch và nâng cấp, có thể nhanh chóng cập nhật phiên bản lõi, hỗ trợ các tính năng Xray-core mới.
* Hỗ trợ UDP và nhiều chức năng khác

## Đặc sắc

| Đặc sắc | v2ray | trojan | tất bóng |
| ------------------------------------------- | ----- | ------ | ----------- |
| Nhận thông tin nút | √ | √ | √ |
| Lấy thông tin người dùng | √ | √ | √ |
| Thống kê lưu lượng người dùng | √ | √ | √ |
| Báo cáo thông tin máy chủ | √ | √ | √ |
| Đăng ký tự động chứng chỉ TLS | √ | √ | √ |
| tự động gia hạn chứng chỉ tls | √ | √ | √ |
| Số người trực tuyến | √ | √ | √ |
| Hạn chế người dùng trực tuyến | √ | √ | √ |
| Quy tắc kiểm toán | √ | √ | √ |
| Giới hạn tốc độ cổng nút | √ | √ | √ |
| Giới hạn tốc độ người dùng | √ | √ | √ |
| DNS tùy chỉnh | √ | √ | √ |
## Hỗ trợ giao diện người dùng

| Panel                                                  | v2ray | trojan | shadowsocks                                 |
| ------------------------------------------------------ | ----- | ------ | ------------------------------------------- |
| [sspanel-uim](https://github.com/Anankke/SSPanel-Uim)  | √     | √      | √ (Đa người dùng một cổng và V2ray-Plugin) |
| [v2board](https://github.com/v2board/v2board)          | √     | √      | √                                           |
| [PMPanel](https://github.com/ByteInternetHK/PMPanel)   | √     | √      | √                                           |
| [ProxyPanel](https://github.com/ProxyPanel/ProxyPanel) | √     | √      | √   

## Cài đặt phần mềm - phát hành
``
wget --no-check-certificate -O AikoR.sh https://raw.githubusercontent.com/AikoCute-Offical/AikoR-Install/master/AikoR.sh && bash AikoR.sh
``
### Một bản cài đặt chính - docker
``
docker pull aikocute / aikor: mới nhất && docker run --restart = always --name aikor -d -v $ {PATCH_TO_CONFIG} /aiko.json:/etc/AikoR/aiko.json --network = host aikocute / aikor: mới nhất
``
### Tệp cấu hình và hướng dẫn chi tiết
Đến sớm
## Telgram

Đến sớm

## Stargazers theo thời gian