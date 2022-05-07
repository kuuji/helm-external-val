#!/usr/bin/env sh

PROJECT_NAME="helm-external-val"
PROJECT_GH="kuuji/$PROJECT_NAME"

# identify system os and arch
OS=$(uname -s)
ARCH=$(uname -m)

if command -v cygpath >/dev/null 2>&1; then
  HELM_BIN="$(cygpath -u "${HELM_BIN}")"
  HELM_PLUGIN_DIR="$(cygpath -u "${HELM_PLUGIN_DIR}")"
fi

[ -z "$HELM_BIN" ] && HELM_BIN=$(command -v helm)

[ -z "$HELM_HOME" ] && HELM_HOME=$(helm env | grep 'HELM_DATA_HOME' | cut -d '=' -f2 | tr -d '"')

mkdir -p "$HELM_HOME"

: "${HELM_PLUGIN_DIR:="$HELM_HOME/plugins/${PROJECT_NAME}"}"

# which mode is the common installer script running in
SCRIPT_MODE="install"
if [ "$1" = "-u" ]; then
  SCRIPT_MODE="update"
fi

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  supported="Darwin-x86_64\nLinux-arm64\nLinux-armv6\nLinux-i386\nLinux-x86_64\nWindows-arm64\nWindows-armv6\nWindows-i386\nWindows-x86_64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuild binary for ${OS}-${ARCH}."
    exit 1
  fi
}

getFormat() {
  if [ "$OS" = "Windows" ]; then
    EXT="zip"
  else
    EXT="tar.gz"
  fi
}

getDownloadURL() {
  version=$(git -C "$HELM_PLUGIN_DIR" describe --tags --exact-match 2>/dev/null || :)
  if [ "$SCRIPT_MODE" = "install" ] && [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/${PROJECT_GH}/releases/download/${version}/${PROJECT_NAME}_${OS}-${ARCH}.${EXT}"
  else
    DOWNLOAD_URL="https://github.com/${PROJECT_GH}/releases/latest/download/${PROJECT_NAME}_${OS}-${ARCH}.${EXT}"
  fi
}

# Temporary dir
mkTempDir() {
  HELM_TMP="$(mktemp -d -t "${PROJECT_NAME}-XXXXXX")"
}

rmTempDir() {
  if [ -d "${HELM_TMP:-/tmp/${PROJECT_NAME}}" ]; then
    rm -rf "${HELM_TMP:-/tmp/${PROJECT_NAME}}"
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="${HELM_TMP}/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if
    command -v curl >/dev/null 2>&1
  then
    curl -sSf -L "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  elif
    command -v wget >/dev/null 2>&1
  then
    wget -q -O - "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  tar xzf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  HELM_TMP_BIN="$HELM_TMP/${PROJECT_NAME}"
  echo "Preparing to install into ${HELM_PLUGIN_DIR}"
  mkdir -p "$HELM_PLUGIN_DIR/bin"
  cp "$HELM_TMP_BIN" "$HELM_PLUGIN_DIR/bin"
}

# exit_trap is executed if on exit (error or not).
exit_trap() {
  result=$?
  rmTempDir
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    printf "\tFor support, go to https://github.com/${PROJECT_GH}.\n"
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_DIR/$PROJECT_NAME"
  "${HELM_PLUGIN_DIR}/bin/${PROJECT_NAME}" -h
  set -e
}

eval `helm env`

# Execution

#Stop execution on any error
trap "exit_trap" EXIT
set -e
verifySupported
getFormat
getDownloadURL
mkTempDir
downloadFile
installFile
testVersion
