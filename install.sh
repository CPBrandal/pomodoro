#!/bin/bash

set -e

echo "Installing Pomodoro Timer..."

REPO="CPBrandal/pomodoro"
VERSION="v2.0.0"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin|linux) ;;
    *)  echo "Unsupported OS: $OS"; exit 1 ;;
esac

BINARY_NAME="pomodoro-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

echo "Detected: ${OS}/${ARCH}"
echo "Downloading ${BINARY_NAME}..."

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary
if command -v curl &> /dev/null; then
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/pomodoro"
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/pomodoro"
else
    echo "Error: curl or wget is required"
    exit 1
fi

chmod +x "$TMP_DIR/pomodoro"

# Install to /usr/local/bin
INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/pomodoro" "$INSTALL_DIR/"
else
    echo "Installing to $INSTALL_DIR requires sudo..."
    sudo mv "$TMP_DIR/pomodoro" "$INSTALL_DIR/"
fi

echo ""
echo "Installation complete! You can now run 'pomodoro' from anywhere."