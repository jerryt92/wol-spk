# JSON Mode

English | [中文](JSON模式.md)

JSON mode, import, and export use the same device shape.

## Recommended format

```json
[
  {
    "name": "Living Room PC",
    "mac": "AA:BB:CC:DD:EE:FF",
    "broadcast": "255.255.255.255",
    "port": 9,
    "notes": "Optional note"
  }
]
```

## Compatible format

Import also accepts an outer object:

```json
{
  "devices": [
    {
      "name": "Living Room PC",
      "mac": "AA:BB:CC:DD:EE:FF",
      "broadcast": "192.168.1.255",
      "port": 9
    }
  ]
}
```

## Fields

| Field | Required | Description |
|-------|----------|-------------|
| `id` | No | Device ID. Generated automatically when empty. |
| `name` | Yes | Device name. |
| `mac` | Yes | MAC address. Supports `AA:BB:CC:DD:EE:FF` and `AA-BB-CC-DD-EE-FF`. |
| `broadcast` | No | Broadcast address. Defaults to `255.255.255.255`. |
| `port` | No | UDP port. Defaults to `9`. |
| `notes` | No | Optional note. |
| `updatedAt` | No | Updated by the app after import/save. |

## Import behavior

- Import overwrites the current device list.
- The UI shows a confirmation prompt before saving.
- MAC addresses are normalized to uppercase colon format.
- Empty `broadcast` becomes `255.255.255.255`.
- Empty `port` becomes `9`.
- Duplicate `id` values are corrected automatically to avoid action conflicts.
