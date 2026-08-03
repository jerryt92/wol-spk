# WOL Manager Test Scenarios

## Test Environment

Record the actual environment before submission:

| Item | Value |
|------|-------|
| NAS model | TBD |
| DSM version | TBD |
| CPU architecture | x86_64 |
| Package version | 1.0.0 |
| Browser | TBD |
| Network | Same LAN |

## 1. Package Build

Steps:

1. Run `./build.sh`.
2. Confirm the SPK file is generated.

Expected result:

- `build/WOLManager-1.0.0-x86_64.spk` exists.
- The SPK contains `INFO`, `package.tgz`, `scripts/`, `conf/`, `PACKAGE_ICON.PNG`, and `PACKAGE_ICON_256.PNG`.

## 2. Manual Installation

Steps:

1. Open DSM Package Center.
2. Click Manual Install.
3. Select `WOLManager-1.0.0-x86_64.spk`.
4. Complete installation.

Expected result:

- Installation succeeds.
- WOL Manager appears in Package Center.
- WOL Manager appears in the DSM main menu.

## 3. Open DSM Application

Steps:

1. Open the DSM main menu.
2. Click WOL Manager.

Expected result:

- WOL Manager opens in a DSM application window.
- The app displays the device list and add device form.
- No browser console or CGI error is shown to the user.

## 4. Add Device

Steps:

1. Enter a device name.
2. Enter a valid MAC address.
3. Keep broadcast address as `255.255.255.255`.
4. Keep port as `9`.
5. Click Save.

Expected result:

- Device is saved successfully.
- Device appears in the device list.
- Device data is persisted after closing and reopening the app.

## 5. Edit Device

Steps:

1. Click Edit on an existing device.
2. Change the name or notes.
3. Click Save.

Expected result:

- Updated information appears in the list.
- Updated information remains after app reload.

## 6. Delete Device

Steps:

1. Click Delete on an existing device.
2. Confirm deletion.

Expected result:

- Device is removed from the list.
- Device remains removed after app reload.

## 7. Wake Device

Steps:

1. Add a target device with Wake-on-LAN enabled.
2. Confirm the target device is powered off or sleeping.
3. Click Wake.

Expected result:

- WOL Manager reports success.
- The NAS sends a Wake-on-LAN magic packet.
- The target device wakes if its BIOS/UEFI, operating system, and network are configured correctly.

Notes:

- If the target does not wake, verify Wake-on-LAN support and network broadcast behavior.
- Test with both `255.255.255.255` and the subnet broadcast address when needed.

## 8. JSON Mode Format

Steps:

1. Switch to JSON Mode.
2. Edit the JSON device list.
3. Click Format.

Expected result:

- Valid JSON is formatted.
- Invalid JSON shows a validation error.

## 9. JSON Mode Save

Steps:

1. Switch to JSON Mode.
2. Replace the device list with a valid JSON array.
3. Click Save JSON.
4. Confirm overwrite.

Expected result:

- Current device list is replaced.
- New devices appear in UI mode.
- Data remains after app reload.

## 10. Export Devices

Steps:

1. Add one or more devices.
2. Click Export.

Expected result:

- Browser downloads `wolmanager-devices.json`.
- Exported JSON contains the current device list.

## 11. Import Devices

Steps:

1. Prepare a valid JSON device list.
2. Click Import.
3. Select the JSON file.
4. Confirm overwrite.

Expected result:

- Current device list is replaced with the imported list.
- Imported devices can be edited, deleted, and woken.

## 12. Invalid Input Validation

Steps:

1. Try saving an empty device name.
2. Try saving an invalid MAC address.
3. Try saving an invalid UDP port outside `1` to `65535`.

Expected result:

- Invalid data is rejected.
- A clear error is displayed.
- Existing saved devices are not corrupted.

## 13. Uninstall

Steps:

1. Export the device list.
2. Uninstall WOL Manager from Package Center.

Expected result:

- Package is removed successfully.
- DSM main menu entry is removed.
- Package data file is removed according to uninstall behavior.

## 14. Reinstall

Steps:

1. Install the SPK again.
2. Open WOL Manager.
3. Import the previously exported JSON file.

Expected result:

- Reinstallation succeeds.
- WOL Manager opens correctly.
- Imported devices appear and can be used.

