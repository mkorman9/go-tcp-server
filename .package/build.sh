#!/usr/bin/env bash

set -e

SCRIPTPATH="$(cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P)"

BUILD_PATH="${BUILD_PATH:-/tmp/go-tcp-server}"
BIN_PATH="${BIN_PATH:-${SCRIPTPATH}/../go-tcp-server}"
OUT_PATH="go-tcp-server_1.0_all.deb"

rm -rf "$BUILD_PATH" || true
cp -a "${SCRIPTPATH}/go-tcp-server" "$BUILD_PATH"
cp "$BIN_PATH" "${BUILD_PATH}/usr/local/bin"
chmod -R 755 "$BUILD_PATH"

dpkg-deb --build --root-owner-group "$BUILD_PATH" "$OUT_PATH"
