# 开发与打包

[English](README.en.md) | 中文

本文档说明 WOL Manager 的代码结构、本地开发、测试和 SPK 打包流程。

## 项目结构

```text
.
├── cmd/wolmanager/        # Go Web/CGI 入口和页面
├── internal/store/        # 设备 JSON 存储和校验
├── internal/wol/          # WOL magic packet 生成和发送
├── synology/              # DSM SPK 元数据、脚本和 UI 入口
├── tools/icongen/         # DSM 图标生成工具
├── build.sh               # 本地打包脚本
└── docs/                  # 文档
```

## 后端 API

| action | 方法 | 用途 |
|--------|------|------|
| `list` | GET | 获取设备列表。 |
| `save` | POST | 新增或更新单个设备。 |
| `delete` | POST | 删除设备。 |
| `replace` | POST | 覆盖整个设备列表，用于 JSON 保存和导入。 |
| `wake` | POST | 唤醒指定设备。 |

## 本地开发

```sh
WOLMANAGER_DATA_DIR=/tmp/wolmanager-dev go run ./cmd/wolmanager
```

打开：

```text
http://127.0.0.1:8088
```

## 图标生成

DSM 应用入口使用高分辨率 `icon_256.png`，并同时生成 `icon_{0}.png` 多尺寸图标作为兼容资源：

```sh
go run ./tools/icongen
```

## 测试

如果 Go 默认缓存目录不可写，可以指定 `GOCACHE`：

```sh
GOCACHE=/private/tmp/wol-spk-gocache go test ./...
```

## 打包

```sh
GOCACHE=/private/tmp/wol-spk-gocache ./build.sh
```

输出：

```text
build/WOLManager-1.0.0-x86_64.spk
```

## GitHub Release

项目内置 GitHub Actions 发布流水线：

```text
.github/workflows/release.yml
```

触发方式：

- 推送 `v*` 标签，例如 `v1.0.0`
- 在 GitHub Actions 页面手动运行 `Release`，输入 `v1.0.0`

发布前建议先提交所有代码，然后创建并推送标签：

```sh
git tag v1.0.0
git push origin v1.0.0
```

流水线会自动执行：

- 安装 Go
- 运行 `go test ./...`
- 执行 `./build.sh`
- 校验 `build/WOLManager-1.0.0-x86_64.spk`
- 校验 Release 标签必须等于 `v` + `synology/INFO` 版本号
- 创建 GitHub Release
- 上传 SPK 到 Release 附件

版本号和包名从 `synology/INFO` 读取，目前保持：

```text
version="1.0.0"
maintainer="jerryt92"
```

SPK 顶层包含：

```text
INFO
conf/privilege
scripts/start-stop-status
package.tgz
```

`package.tgz` 包含：

```text
bin/wolmanager
ui/config
ui/index.cgi
ui/images/icon_256.png
ui/images/icon_{0}.png
```
