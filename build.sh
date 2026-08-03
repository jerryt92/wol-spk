#!/usr/bin/env bash
set -euo pipefail

PACKAGE="WOLManager"
VERSION="1.0.0"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${ROOT_DIR}/build"
STAGE_DIR="${BUILD_DIR}/stage"
TARGET_DIR="${STAGE_DIR}/target"
SPK_DIR="${BUILD_DIR}/spk"

rm -rf "${BUILD_DIR}"
mkdir -p "${TARGET_DIR}/bin" "${TARGET_DIR}/ui" "${SPK_DIR}"

echo "Building Go binary..."
go run ./tools/i18ngen
go run ./tools/icongen
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "${TARGET_DIR}/bin/wolmanager" ./cmd/wolmanager

echo "Preparing DSM package payload..."
cp -R "${ROOT_DIR}/synology/ui/." "${TARGET_DIR}/ui/"
chmod 755 "${TARGET_DIR}/ui/index.cgi"
chmod 755 "${TARGET_DIR}/bin/wolmanager"

tar -C "${TARGET_DIR}" -czf "${SPK_DIR}/package.tgz" .
cp "${ROOT_DIR}/synology/PACKAGE_ICON.PNG" "${SPK_DIR}/PACKAGE_ICON.PNG"
cp "${ROOT_DIR}/synology/PACKAGE_ICON_256.PNG" "${SPK_DIR}/PACKAGE_ICON_256.PNG"
cp "${ROOT_DIR}/synology/INFO" "${SPK_DIR}/INFO"
mkdir -p "${SPK_DIR}/conf" "${SPK_DIR}/scripts"
cp "${ROOT_DIR}/synology/conf/privilege" "${SPK_DIR}/conf/privilege"
cp -R "${ROOT_DIR}/synology/scripts/." "${SPK_DIR}/scripts/"
chmod 755 "${SPK_DIR}/scripts/"*

SPK_PATH="${BUILD_DIR}/${PACKAGE}-${VERSION}-x86_64.spk"
tar -C "${SPK_DIR}" -cf "${SPK_PATH}" INFO PACKAGE_ICON.PNG PACKAGE_ICON_256.PNG conf scripts package.tgz

echo "Created ${SPK_PATH}"
