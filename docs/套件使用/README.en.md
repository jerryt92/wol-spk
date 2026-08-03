# Package Usage

English | [中文](README.md)

This document covers daily WOL Manager usage.

## Open the app

After installation, open **WOL Manager** from the DSM main menu.

If the DSM main menu does not show the entry yet, refresh DSM or sign out and sign in again. If it is still missing, see [Troubleshooting](../故障排查/README.en.md).

## Add a device

1. Click Add Device.
2. Enter a device name.
3. Enter a MAC address, for example `AA:BB:CC:DD:EE:FF`.
4. Leave broadcast address as `255.255.255.255` unless your network needs another value.
5. Leave port as `9` unless your device expects another WOL port.
6. Click Save.

If your network blocks global broadcast, use the subnet broadcast address instead, such as `192.168.1.255`.

## Wake a device

Click Wake on a device card. The app sends a Wake-on-LAN magic packet.

The target device must have Wake-on-LAN enabled in BIOS/UEFI, network adapter settings, or the operating system.

## Daily maintenance

Use UI mode for small lists. Use JSON mode, import, and export from [Device Management](../设备管理/README.en.md) when you need bulk editing, backups, or migration.
