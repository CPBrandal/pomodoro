#!/bin/bash

set -e

echo "Installing Pomodoro Timer..."

# Build the binary
echo "Building..."
go build -o pomodoro main.go

# Determine install location
INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
    mv pomodoro "$INSTALL_DIR/"
    echo "Installed to $INSTALL_DIR/pomodoro"
else
    echo "Installing to $INSTALL_DIR requires sudo..."
    sudo mv pomodoro "$INSTALL_DIR/"
    echo "Installed to $INSTALL_DIR/pomodoro"
fi

echo ""
echo "Installation complete! You can now run 'pomodoro' from anywhere."