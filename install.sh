#!/usr/bin/env bash
set -euo pipefail

BINARY="razorpay"
INSTALL_DIR="/usr/local/bin"
S3_BASE_URL="https://rzp-1018-nonprod-betasite.razorpay.com.s3.ap-south-1.amazonaws.com/cli"

# ---------------------------------------------------------------------------
# Detect OS and architecture
# ---------------------------------------------------------------------------
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin) OS_NAME="mac-os" ;;
  Linux)  OS_NAME="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)        ARCH_NAME="x86_64" ;;
  arm64|aarch64) ARCH_NAME="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# ---------------------------------------------------------------------------
# Resolve latest version from S3
# ---------------------------------------------------------------------------
echo "Fetching latest version..."
VERSION=$(curl -fsSL "${S3_BASE_URL}/latest")

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version. Check your internet connection."
  exit 1
fi

echo "Latest version: $VERSION"

# ---------------------------------------------------------------------------
# Download and extract
# ---------------------------------------------------------------------------
ARCHIVE="${BINARY}_${VERSION#v}_${OS_NAME}_${ARCH_NAME}.tar.gz"
URL="${S3_BASE_URL}/${VERSION#v}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"

echo "Downloading $ARCHIVE..."
curl -fsSL "$URL" -o "$TMP_DIR/$ARCHIVE"

echo "Extracting..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

# ---------------------------------------------------------------------------
# Install
# ---------------------------------------------------------------------------
if [ ! -w "$INSTALL_DIR" ]; then
  echo "Installing to $INSTALL_DIR (requires sudo)..."
  sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
  sudo chmod +x "$INSTALL_DIR/$BINARY"
else
  mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
  chmod +x "$INSTALL_DIR/$BINARY"
fi

rm -rf "$TMP_DIR"

echo ""
echo "razorpay $VERSION installed to $INSTALL_DIR/$BINARY"
echo "Run 'razorpay configure' to set up your API credentials."
