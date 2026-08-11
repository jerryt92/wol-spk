# Development and Packaging

English | [中文](README.md)

This document covers WOL Manager code layout, local development, tests, and SPK packaging.

## Project layout

```text
.
├── cmd/wolmanager/        # Go Web/CGI entrypoint and page
├── internal/store/        # Device JSON storage and validation
├── internal/wol/          # WOL magic packet generation and sending
├── synology/              # DSM SPK metadata, scripts, and UI entry
├── tools/icongen/         # DSM icon generator
├── build.sh               # macOS/Linux package script
├── build.ps1              # Windows PowerShell package script
└── docs/                  # Documentation
```

## Backend API

| action | Method | Purpose |
|--------|--------|---------|
| `list` | GET | List devices. |
| `save` | POST | Add or update one device. |
| `delete` | POST | Delete a device. |
| `replace` | POST | Replace the full device list for JSON save/import. |
| `wake` | POST | Wake a device. |

## Local development

```sh
WOLMANAGER_DATA_DIR=/tmp/wolmanager-dev go run ./cmd/wolmanager
```

Open:

```text
http://127.0.0.1:8088
```

## Icon generation

The DSM app entry uses the high-resolution `icon_256.png` icon and also generates multi-size `icon_{0}.png` icons as compatibility assets:

```sh
go run ./tools/icongen
```

## Tests

If the default Go cache directory is not writable, set `GOCACHE`:

```sh
GOCACHE=/private/tmp/wol-spk-gocache go test ./...
```

## Package

```sh
GOCACHE=/private/tmp/wol-spk-gocache ./build.sh
```

Windows PowerShell:

```powershell
./build.ps1
```

Output:

```text
build/<package>-<version>-x86_64.spk
```

## GitHub Release

The project includes a GitHub Actions release workflow:

```text
.github/workflows/release.yml
```

Triggers:

- Push a `v*` tag that matches the version in `synology/INFO`.
- Run `Release` manually from the GitHub Actions page and enter the matching tag.

Commit all changes first, then create and push a tag:

```sh
git tag v<version>
git push origin v<version>
```

The workflow automatically:

- Sets up Go
- Runs `go test ./...`
- Runs `./build.sh`
- Verifies `build/<package>-<version>-x86_64.spk`
- Verifies that the release tag equals `v` plus the `synology/INFO` version
- Creates a GitHub Release
- Uploads the SPK as a release asset

The package name and version are maintained only in `synology/INFO`:

```text
version="<version>"
maintainer="jerryt92"
```

SPK top-level files:

```text
INFO
conf/privilege
scripts/start-stop-status
package.tgz
```

`package.tgz` contains:

```text
bin/wolmanager
ui/config
ui/index.cgi
ui/images/icon_256.png
ui/images/icon_{0}.png
```
