#!/bin/bash

set -e

echo "Building Pomodoro Timer for all platforms..."

# Output directory
OUT_DIR="dist"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

# Build targets
TARGETS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

for TARGET in "${TARGETS[@]}"; do
    OS="${TARGET%/*}"
    ARCH="${TARGET#*/}"
    OUTPUT_NAME="pomodoro-${OS}-${ARCH}"
    
    echo "Building ${OUTPUT_NAME}..."
    GOOS=$OS GOARCH=$ARCH go build -o "$OUT_DIR/$OUTPUT_NAME" .
done

echo ""
echo "Build complete! Binaries are in the '$OUT_DIR' directory:"
ls -la "$OUT_DIR"