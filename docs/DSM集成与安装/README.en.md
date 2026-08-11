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
4. Select `build/<package>-<version>-x86_64.spk` (the package name and version come from `synology/INFO`).
5. Open WOL Manager from the DSM main menu after installation.

## Checking GitHub updates

The **Check for Updates** button in the upper-right corner reads the latest stable [GitHub Release](https://github.com/jerryt92/wol-spk/releases/latest) and compares it with the installed version. When an update is available, confirm the download to retrieve the SPK matching the NAS architecture, then install that file manually through DSM Package Center.

Each release must meet both requirements:

- Its tag is `v<version>`, for example `v1.0.3`.
- Its assets include `WOLManager-<version>-x86_64.spk`.

The app only checks for and downloads updates; it never replaces a running package itself.

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

During an upgrade, the package copies this file to DSM's temporary upgrade workspace and restores it after the new version is installed, so normal upgrades preserve devices. The file is removed only on uninstall.

The uninstall script removes `devices.json`. Export the device list from the app first if you want to keep a backup.

## Why not synonet

DSM 7 packages normally run as a low-privilege package user, while `synonet --wake` may require root on some systems. WOL Manager sends UDP magic packets directly from Go to avoid root-only commands and extra script setup.
