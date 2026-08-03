# WOL Manager - Synology Package Center Submission

## Product Information

Product name: WOL Manager

Product category: Utilities

Package ID: WOLManager

Version: 1.0.0

Maintainer: jerryt92

Website: https://jerryt92.top/

Support email: jerrytian92@outlook.com

License / pricing: Free

## Short Description

WOL Manager is a Wake-on-LAN utility for Synology DSM. It allows users to manage local network devices and send Wake-on-LAN magic packets directly from the NAS.

## Full Description

WOL Manager is a native DSM 7 package for managing Wake-on-LAN devices on a Synology NAS. Users can add devices with a name, MAC address, broadcast address, port, and optional notes, then wake those devices from the DSM interface.

The package also includes JSON mode for bulk editing, import, export, backup, and migration of the device list. Device data is stored locally on the NAS at `/var/packages/WOLManager/var/devices.json`.

WOL Manager is built as a native SPK package and does not require Docker or an additional runtime such as Python, Java, Node.js, or npm on the NAS. The backend is a static Go binary and the DSM UI is integrated through the DSM main menu.

## Key Features

- Add, edit, delete, and wake Wake-on-LAN devices.
- Store device name, MAC address, broadcast address, UDP port, and notes.
- Send UDP Wake-on-LAN magic packets from the NAS.
- JSON mode for bulk device management.
- Import and export device lists as JSON.
- Native DSM 7 application window integration.
- No Docker or external runtime dependency on the NAS.
- Runs as the package user according to DSM package privilege settings.

## Compatibility

| Item | Value |
|------|-------|
| DSM version | DSM 7.0 or later |
| CPU architecture | x86_64 |
| Initial target devices | Intel/AMD 64-bit Synology NAS models |
| Runtime dependency | None on NAS |
| Network requirement | Local network UDP broadcast or subnet broadcast |

## Package Behavior

### Installation

The package installs a native DSM application and registers a DSM main menu entry named WOL Manager.

### Runtime

WOL Manager is implemented as a CGI-style DSM UI application. It does not run a background daemon. User actions trigger local CGI requests that read/write the device list or send Wake-on-LAN UDP packets.

### Data Storage

Device data is stored locally:

```text
/var/packages/WOLManager/var/devices.json
```

### Uninstall

The uninstall process removes the package data file. Users should export the device list before uninstalling if they want to keep a backup.

## Security and Permissions

- The package runs as the DSM package user.
- It does not require root privileges.
- It does not use `synonet --wake`.
- It sends UDP Wake-on-LAN packets directly from the Go backend.
- Device data is stored locally on the NAS.
- There is no external cloud service dependency.

## Review Materials Checklist

- [ ] SPK file: `build/WOLManager-1.0.0-x86_64.spk`
- [x] Package icon: `synology/PACKAGE_ICON.PNG`
- [x] Package icon 256: `synology/PACKAGE_ICON_256.PNG`
- [x] Screenshot: `release-materials/screenshots/01-device-list-ui-mode.png`
- [x] Screenshot: `release-materials/screenshots/02-json-mode.png`
- [x] Operation manual: `release-materials/operation-manual.md`
- [x] Test scenarios: `release-materials/test-scenarios.md`
- [x] Changelog: `release-materials/changelog.md`
- [ ] Public download link for SPK

## Screenshots

### Device List and UI Mode

File: `release-materials/screenshots/01-device-list-ui-mode.png`

This screenshot shows the main DSM application window, device list, Wake button, import/export actions, and add device form.

### JSON Mode

File: `release-materials/screenshots/02-json-mode.png`

This screenshot shows JSON mode for bulk editing, formatting, reload, import, and export workflows.

