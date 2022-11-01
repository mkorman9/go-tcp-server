#!/usr/bin/env bash

set -e

DEB_FILE="go-tcp-server_1.0.0_amd64.deb"
REMOTE_HOST="localhost"
REMOTE_USER="root"
REMOTE_DEB_PATH="."
SUDO=""
KEY_OPTS=""

while [[ $# -gt 0 ]]; do
  case $1 in
    -f|--file)
      DEB_FILE="$2"
      shift
      shift
      ;;
    -h|--host)
      REMOTE_HOST="$2"
      shift
      shift
      ;;
    -u|--user)
      REMOTE_USER="$2"
      shift
      shift
      ;;
    -p|--path)
      REMOTE_DEB_PATH="$2"
      shift
      shift
      ;;
    -k|--key)
      KEY_OPTS="-i $2"
      shift
      shift
      ;;
    -*|--*)
      echo "Unknown option $1"
      echo "usage: ./deploy.sh [OPTIONS]"
      echo "OPTIONS:"
      echo "  --file / -f <name of .deb file to deploy> (default: go-tcp-server_1.0.0_amd64.deb)"
      echo "  --host / -h <remote host to deploy to> (default: localhost)"
      echo "  --user / -u <remote user to log in to> (default: root)"
      echo "  --path / -p <path on a remote host to upload .deb file> (default: .)"
      echo "  --key / -k <path to SSH private key>"
      exit 1
      ;;
    *)
      shift
      ;;
  esac
done

if [[ ! -f "$DEB_FILE" ]]; then
  echo "DEB file couldn't be found: ${DEB_FILE}"
  exit 1
fi

if [[ "$REMOTE_USER" != "root" ]]; then
  SUDO="sudo "
fi

SSH_OPTS="-o BatchMode=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR $KEY_OPTS"

echo "Checking remote host..." && \
  ssh ${SSH_OPTS} "${REMOTE_USER}@${REMOTE_HOST}" -C "uname -a" && \
  echo "Uploading package..." && \
  scp ${SSH_OPTS} "${DEB_FILE}" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DEB_PATH}" && \
  echo "Starting installation..." && \
  ssh ${SSH_OPTS} "${REMOTE_USER}@${REMOTE_HOST}" -C "${SUDO}dpkg -i ${REMOTE_DEB_PATH}/${DEB_FILE}"
