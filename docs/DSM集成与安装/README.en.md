# DSM Integration and Installation

English | [中文](README.md)

This document explains how WOL Manager integrates with DSM and how to install it.

## Package type

WOL Manager is a native DSM 7 SPK package:

- No Docker.
- No Python, Java, Node, or npm required on the NAS.
- The backend is a single Linux amd64 binary built with Go.
- DSM links the UI directory to `/webman/3rdparty/WOLManager/` through `dsmuidir="ui"`.

## Compatibility

| Item | Value |
|------|-------|
| DSM | 7.0 or later |
| Architecture | `x86_64` |
| Initial target | Intel/AMD 64-bit Synology devices such as DS225+ |

## Installation

1. Run `./build.sh` on the build machine.
2. Open DSM Package Center.
3. Choose Manual Install.
4. Select `build/WOLManager-1.0.0-x86_64.spk`.
5. Open WOL Manager from the DSM main menu after installation.

## DSM desktop entry

The package entry is configured in `synology/ui/config` using DSM `app` mode. Clicking the main menu icon creates an internal DSM application window, and `wolmanager.js` embeds `/webman/3rdparty/WOLManager/index.cgi` inside that window.

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

If the main menu entry is missing after installation, refresh DSM or sign out and sign in again. See [Troubleshooting](../故障排查/README.en.md) for more checks.

## Package Center icon

The SPK top level contains `PACKAGE_ICON.PNG` and `PACKAGE_ICON_256.PNG` for DSM Package Center. After installation, `postinst` also copies the same icons to `/var/packages/WOLManager/`, keeping Package Center, DSM desktop, and the app UI visually consistent.

## Data path

Device data is stored at:

```text
/var/packages/WOLManager/var/devices.json
```

The uninstall script removes `devices.json`. Export the device list from the app first if you want to keep a backup.

## Why not synonet

DSM 7 packages normally run as a low-privilege package user, while `synonet --wake` may require root on some systems. WOL Manager sends UDP magic packets directly from Go to avoid root-only commands and extra script setup.
