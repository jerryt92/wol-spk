# WOL Manager Changelog

## 1.0.0

Initial release.

### Added

- Native DSM 7 SPK package for Wake-on-LAN device management.
- DSM main menu application integration.
- Device list UI for adding, editing, deleting, and waking devices.
- Wake-on-LAN magic packet sending from the NAS.
- MAC address validation and normalized input.
- Configurable broadcast address and UDP port.
- Optional device notes.
- JSON mode for bulk editing.
- JSON import and export for backup and migration.
- Local data storage under `/var/packages/WOLManager/var/devices.json`.
- Package icons for DSM Package Center and DSM application UI.

### Compatibility

- DSM 7.0 or later.
- x86_64 Synology NAS devices.

### Notes

- WOL Manager does not require Docker or an additional runtime on the NAS.
- The package runs as the DSM package user.
- Wake-on-LAN success depends on target device BIOS/UEFI, operating system, network adapter, and network broadcast configuration.

