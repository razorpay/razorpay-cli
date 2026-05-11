#!/usr/bin/env bash
set -euo pipefail

BINARY="razorpay"
INSTALL_DIR="$HOME/.local/bin"
BASE_URL="https://razorpay.com/cli/latest"

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
# Download and extract
# ---------------------------------------------------------------------------
ARCHIVE="${BINARY}_${OS_NAME}_${ARCH_NAME}.tar.gz"
URL="${BASE_URL}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"

echo "Downloading ${ARCHIVE}..."
curl -fsSL "$URL" -o "$TMP_DIR/$ARCHIVE"

echo "Extracting..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

# ---------------------------------------------------------------------------
# Install
# ---------------------------------------------------------------------------
mkdir -p "$INSTALL_DIR"
mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

rm -rf "$TMP_DIR"

# ---------------------------------------------------------------------------
# Ensure ~/.local/bin is in PATH
# ---------------------------------------------------------------------------
if ! echo "$PATH" | tr ':' '\n' | grep -q "^$INSTALL_DIR$"; then
  SHELL_NAME="$(basename "$SHELL")"
  case "$SHELL_NAME" in
    zsh)  PROFILE="$HOME/.zshrc" ;;
    bash) PROFILE="$HOME/.bashrc" ;;
    *)    PROFILE="$HOME/.profile" ;;
  esac

  EXPORT_LINE="export PATH=\"\$HOME/.local/bin:\$PATH\""

  if ! grep -qF '.local/bin' "$PROFILE" 2>/dev/null; then
    echo "" >> "$PROFILE"
    echo "$EXPORT_LINE" >> "$PROFILE"
    echo "Added $INSTALL_DIR to PATH in $PROFILE"
  fi

  echo ""
  echo "NOTE: Run 'source $PROFILE' or open a new terminal for the PATH change to take effect."
fi

VERSION=$("$INSTALL_DIR/$BINARY" --version 2>/dev/null || echo "unknown")
echo ""
echo "razorpay ${VERSION} installed to $INSTALL_DIR/$BINARY"
echo "Run 'razorpay configure' to set up your API credentials."
