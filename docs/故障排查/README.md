# 故障排查

[English](README.en.md) | 中文

本文档整理 WOL Manager 常见问题。

## 设备没有被唤醒

检查目标设备：

- BIOS/UEFI 已开启 Wake-on-LAN。
- 网卡驱动允许关机后唤醒。
- Windows 中关闭「快速启动」后再测试。
- 目标设备使用有线网卡时通常更稳定。

检查网络：

- NAS 和目标设备在同一个局域网。
- 如果 `255.255.255.255` 不工作，尝试网段广播地址，例如 `192.168.1.255`。
- 路由器、交换机或 VLAN 没有阻止 UDP 广播。
- 端口通常是 `9`，部分设备或网络习惯使用 `7`。

## DSM 页面打不开或入口不可见

检查套件是否已安装，并确认 DSM 主菜单中有 WOL Manager。

如果安装后没有入口：

- 刷新 DSM 页面。
- 退出 DSM 后重新登录。
- 确认安装包中 `synology/ui/config` 使用 `app` 入口，并且 `synology/ui/wolmanager.js` 存在。
- 直接访问 `/webman/3rdparty/WOLManager/index.cgi` 测试入口是否存在。

## JSON 导入失败

常见原因：

- JSON 不是数组，也不是包含 `devices` 数组的对象。
- 某个设备缺少 `name`。
- 某个设备缺少 `mac`。
- MAC 地址格式不正确。
- `port` 不在 `1` 到 `65535` 之间。

可以先切到 JSON 模式点击「格式化」，定位语法问题。

## 本地构建失败

如果报 Go 缓存目录不可写，使用：

```sh
GOCACHE=/private/tmp/wol-spk-gocache ./build.sh
```
