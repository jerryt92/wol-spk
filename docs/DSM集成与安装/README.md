# DSM 集成与安装

[English](README.en.md) | 中文

本文档说明 WOL Manager 的 DSM 原生集成方式和安装流程。

## 套件定位

WOL Manager 是 DSM 7 原生 SPK 套件：

- 不使用 Docker。
- 不要求 NAS 额外安装 Python、Java、Node 或 npm。
- 后端是 Go 编译出的 Linux amd64 单文件。
- DSM 通过 `dsmuidir="ui"` 将 UI 目录挂载到 `/webman/3rdparty/WOLManager/`。

## 兼容性

| 项目 | 值 |
|------|----|
| DSM | 7.0 及以上 |
| 架构 | `x86_64` |
| 首版目标 | DS225+ 这类 Intel/AMD 64-bit 群晖设备 |

## 安装

1. 在开发机执行 `./build.sh`。
2. 打开 DSM「套件中心」。
3. 选择「手动安装」。
4. 选择 `build/<package>-<version>-x86_64.spk`（包名和版本从 `synology/INFO` 读取）。
5. 安装完成后，从 DSM 主菜单打开 WOL Manager。

## 检查 GitHub 更新

应用右上角的「检查更新」会读取 [GitHub Releases](https://github.com/jerryt92/wol-spk/releases/latest) 中的最新正式发布版本，并与当前已安装版本比较。发现更新后，确认下载即可取得与 NAS 架构匹配的 SPK；再回到 DSM「套件中心」手动安装该文件完成升级。

发布新版本时请同时满足：

- 标签必须是 `v<version>`，例如 `v1.0.3`。
- Release 附件必须包含 `WOLManager-<version>-x86_64.spk`。

应用仅检查和下载更新，不会自行覆盖已运行的套件。

## DSM 桌面入口

套件入口配置位于 `synology/ui/config`，使用 DSM `app` 模式。点击主菜单图标后，DSM 会创建一个内部应用窗口，窗口内容由 `wolmanager.js` 嵌入 `/webman/3rdparty/WOLManager/index.cgi`。

```json
{
  "wolmanager.js": {
    "io.github.jerryt92.spk.wol": {
      "type": "app",
      "appWindow": "io.github.jerryt92.spk.wol.MainWindow"
    }
  }
}
```

如果安装后主菜单没有入口，先刷新 DSM 页面或重新登录。更多检查项见 [故障排查](../故障排查/README.md)。

## 套件中心图标

SPK 顶层包含 `PACKAGE_ICON.PNG` 和 `PACKAGE_ICON_256.PNG`，用于 DSM 套件中心。安装后 `postinst` 也会把同一套图标复制到 `/var/packages/WOLManager/`，确保套件中心、DSM 桌面和应用 UI 使用统一视觉。

## 数据路径

设备列表保存在：

```text
/var/packages/WOLManager/var/devices.json
```

升级时，套件会先将该文件复制到 DSM 的临时升级目录，并在新版本安装后恢复；正常升级不会丢失设备。卸载时才会删除该文件。

卸载套件时，卸载脚本会删除 `devices.json`。如果需要保留设备列表，请先在应用中使用「导出」备份 JSON。

## 为什么不用 synonet

DSM 7 套件通常以 package 低权限用户运行，而 `synonet --wake` 在部分系统上可能需要 root 权限。WOL Manager 直接用 Go 发送 UDP magic packet，避免依赖 root 权限和额外脚本配置。
