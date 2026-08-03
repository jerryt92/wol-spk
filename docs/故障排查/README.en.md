# Troubleshooting

English | [中文](README.md)

This document covers common WOL Manager issues.

## Device does not wake

Check the target device:

- Wake-on-LAN is enabled in BIOS/UEFI.
- The network adapter allows wake after shutdown.
- On Windows, disable Fast Startup and test again.
- Wired Ethernet is usually more reliable for WOL.

Check the network:

- The NAS and target device are on the same LAN.
- If `255.255.255.255` does not work, try the subnet broadcast address, such as `192.168.1.255`.
- The router, switch, or VLAN is not blocking UDP broadcast.
- The usual port is `9`; some devices or networks use `7`.

## DSM page does not open or entry is missing

Check that the package is installed and WOL Manager appears in the DSM main menu.

If the app entry is missing:

- Refresh DSM.
- Sign out and sign in again.
- Confirm `synology/ui/config` uses the `app` entry and `synology/ui/wolmanager.js` exists.
- Open `/webman/3rdparty/WOLManager/index.cgi` directly to test the linked UI.

## JSON import fails

Common causes:

- JSON is not an array and is not an object containing a `devices` array.
- A device is missing `name`.
- A device is missing `mac`.
- A MAC address is invalid.
- `port` is outside `1` to `65535`.

You can switch to JSON mode and click Format to locate syntax issues.

## Local build fails

If the Go cache directory is not writable, use:

```sh
GOCACHE=/private/tmp/wol-spk-gocache ./build.sh
```
