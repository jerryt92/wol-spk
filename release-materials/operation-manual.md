# WOL Manager Operation Manual

## Overview

WOL Manager is a DSM utility for managing Wake-on-LAN devices and sending Wake-on-LAN magic packets from a Synology NAS.

## Open WOL Manager

1. Install the SPK package in DSM Package Center.
2. Open the DSM main menu.
3. Click WOL Manager.

If the main menu entry does not appear immediately, refresh DSM or sign out and sign in again.

## Add a Device

1. Click Add Device.
2. Enter a device name.
3. Enter the MAC address.
4. Keep the default broadcast address `255.255.255.255`, or use a subnet broadcast address such as `192.168.1.255` if required by the network.
5. Keep the default UDP port `9`, or change it if the device/network uses another WOL port.
6. Optionally enter notes.
7. Click Save.

## Wake a Device

1. Confirm the target device supports Wake-on-LAN.
2. Confirm Wake-on-LAN is enabled in BIOS/UEFI, network adapter settings, or the operating system.
3. Open WOL Manager.
4. Click Wake on the target device.

The NAS sends a Wake-on-LAN magic packet to the configured broadcast address and port.

## Edit a Device

1. Click Edit on a device.
2. Update the name, MAC address, broadcast address, port, or notes.
3. Click Save.

## Delete a Device

1. Click Delete on a device.
2. Confirm the deletion prompt.

## JSON Mode

JSON mode is designed for bulk editing, backups, and migration.

1. Click JSON Mode.
2. Edit the JSON device list.
3. Click Format to validate and auto-indent the JSON.
4. Click Save JSON to replace the current device list.

The JSON file can be either:

```json
[
  {
    "name": "Office PC",
    "mac": "AA:BB:CC:DD:EE:FF",
    "broadcast": "255.255.255.255",
    "port": 9,
    "notes": ""
  }
]
```

Or an object with a `devices` array:

```json
{
  "devices": [
    {
      "name": "Office PC",
      "mac": "AA:BB:CC:DD:EE:FF",
      "broadcast": "255.255.255.255",
      "port": 9,
      "notes": ""
    }
  ]
}
```

## Export Devices

1. Click Export.
2. Save the downloaded `wolmanager-devices.json` file.

Export is recommended before uninstalling the package or migrating to another NAS.

## Import Devices

1. Click Import.
2. Select a JSON file.
3. Confirm the overwrite prompt.

The current device list will be replaced with the imported list.

## Troubleshooting

If a device does not wake:

- Confirm Wake-on-LAN is enabled on the target device.
- Use wired Ethernet when possible.
- Confirm the NAS and the target device are on the same LAN or that broadcast routing is configured.
- Try a subnet broadcast address such as `192.168.1.255`.
- Try UDP port `9` or `7`.
- Confirm routers, switches, VLANs, or firewall rules do not block UDP broadcast packets.

If the app does not open:

- Refresh DSM.
- Sign out and sign in again.
- Confirm the package is installed in Package Center.

