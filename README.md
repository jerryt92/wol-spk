# WOL Manager for Synology DSM

DSM 7 native SPK package for waking devices on the local network. No Docker, no extra runtime on the NAS.

## Language / 语言

- [中文文档](docs/README.md)
- [English documentation](docs/README.en.md)

## Documentation

Documents are organized by feature:

- [套件使用](docs/套件使用/README.md) / [Package Usage](docs/套件使用/README.en.md)
- [DSM 集成与安装](docs/DSM集成与安装/README.md) / [DSM Integration and Installation](docs/DSM集成与安装/README.en.md)
- [设备管理](docs/设备管理/README.md) / [Device Management](docs/设备管理/README.en.md)
- [开发与打包](docs/开发与打包/README.md) / [Development and Packaging](docs/开发与打包/README.en.md)
- [故障排查](docs/故障排查/README.md) / [Troubleshooting](docs/故障排查/README.en.md)

## Build

```sh
./build.sh
```

The generated package is:

```text
build/WOLManager-1.0.0-x86_64.spk
```
