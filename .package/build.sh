#!/usr/bin/env bash

set -e

SCRIPTPATH="$(cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P)"

BUILD_PATH="${BUILD_PATH:-/tmp/go-tcp-server}"
BIN_PATH="${BIN_PATH:-${SCRIPTPATH}/../go-tcp-server}"
VERSION="${VERSION:-1.0.0}"
OUT_PATH="go-tcp-server_${VERSION}_amd64.deb"

rm -rf "$BUILD_PATH" || true
cp -a "${SCRIPTPATH}/go-tcp-server" "$BUILD_PATH"
find "$BUILD_PATH" -type f -exec sed -i "s/__VERSION__/${VERSION}/g" {} \;
cp "$BIN_PATH" "${BUILD_PATH}/usr/local/bin"
chmod -R 755 "$BUILD_PATH"

dpkg-deb --build --root-owner-group "$BUILD_PATH" "$OUT_PATH"
