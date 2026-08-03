# Device Management

English | [中文](README.md)

This document covers WOL Manager device list management.

## UI mode

UI mode is for daily maintenance of small device lists:

- Add devices.
- Edit devices.
- Delete devices.
- Wake devices.

The app validates required fields and normalizes MAC addresses on save.

## JSON mode

JSON mode is for bulk editing:

- Click JSON Mode to open the editor.
- Click Format to validate and auto-indent the JSON.
- Click Save JSON to overwrite the current device list. A confirmation prompt is shown first.
- Click Reload to discard unsaved editor changes and restore the current device list.

See [JSON Mode](JSON-Mode.en.md) for details.

## Export JSON

Click Export to download:

```text
wolmanager-devices.json
```

The export is a device array, suitable for backup and migration.

Uninstalling the package removes the device list stored on the NAS. Export JSON first if you want to keep it.

## Import JSON

1. Click Import.
2. Select a `.json` file.
3. Confirm the overwrite prompt.
4. The current device list is replaced with the imported file.

The import file can be a device array or an object containing a `devices` array.
