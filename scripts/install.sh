#!/usr/bin/env bash

PROJECT_NAME="helm-external-val"
PROJECT_GH="kuuji/$PROJECT_NAME"

: "${HELM_PLUGIN_DIR:="$HELM_PLUGINS/$PROJECT_NAME"}"

# which mode is the common installer script running in
SCRIPT_MODE="install"
if [ "$1" = "-u" ]; then
  SCRIPT_MODE="update"
fi

getDownloadURL() {
  version=$(git -C "$HELM_PLUGIN_DIR" describe --tags --exact-match 2>/dev/null || :)
  if [ "$SCRIPT_MODE" = "install" ] && [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/$PROJECT_NAME"
  else
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/latest/download/$PROJECT_NAME"
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  echo "Downloading $DOWNLOAD_URL"
  if
    command -v curl >/dev/null 2>&1
  then
    curl -sSf -L "$DOWNLOAD_URL" >"$HELM_PLUGIN_DIR/$PROJECT_NAME"
  elif
    command -v wget >/dev/null 2>&1
  then
    wget -q -O - "$DOWNLOAD_URL" >"$HELM_PLUGIN_DIR/$PROJECT_NAME"
  fi
}


eval `helm env`

getDownloadURL
downloadFile
chmod +x "$HELM_PLUGIN_DIR/$PROJECT_NAME"
# getDownloadURL checks the latest available version.
