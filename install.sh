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
  Darwin) OS_NAME="mac-os"; CHECKSUM_OS="mac" ;;
  Linux)  OS_NAME="linux";  CHECKSUM_OS="linux" ;;
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
# Download
# ---------------------------------------------------------------------------
ARCHIVE="${BINARY}_${OS_NAME}_${ARCH_NAME}.tar.gz"
URL="${BASE_URL}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${ARCHIVE}..."
curl -fsSL "$URL" -o "$TMP_DIR/$ARCHIVE"

# ---------------------------------------------------------------------------
# Verify checksum
#
# The published checksum file lists the *versioned* archive name, while the
# latest/ prefix serves it without the version -- so resolve the version and
# rebuild the name to look it up.
# ---------------------------------------------------------------------------
echo "Verifying checksum..."

VERSION="$(curl -fsSL "${BASE_URL}/version")"
VERSION_NUM="${VERSION#v}"
CHECKSUM_FILE="${BINARY}-${CHECKSUM_OS}-checksums.txt"
curl -fsSL "${BASE_URL}/${CHECKSUM_FILE}" -o "$TMP_DIR/$CHECKSUM_FILE"

VERSIONED_ARCHIVE="${BINARY}_${VERSION_NUM}_${OS_NAME}_${ARCH_NAME}.tar.gz"
EXPECTED_SHA="$(awk -v name="$VERSIONED_ARCHIVE" '$2 == name { print $1 }' "$TMP_DIR/$CHECKSUM_FILE")"

if [ -z "$EXPECTED_SHA" ]; then
  echo "Error: no checksum published for ${VERSIONED_ARCHIVE} in ${CHECKSUM_FILE}." >&2
  echo "Refusing to install an unverified binary." >&2
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL_SHA="$(sha256sum "$TMP_DIR/$ARCHIVE" | awk '{ print $1 }')"
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL_SHA="$(shasum -a 256 "$TMP_DIR/$ARCHIVE" | awk '{ print $1 }')"
else
  echo "Error: neither sha256sum nor shasum is available; cannot verify the download." >&2
  exit 1
fi

if [ "$EXPECTED_SHA" != "$ACTUAL_SHA" ]; then
  echo "Error: checksum mismatch for ${ARCHIVE}." >&2
  echo "  expected: ${EXPECTED_SHA}" >&2
  echo "  actual:   ${ACTUAL_SHA}" >&2
  echo "Refusing to install. Please report this at https://github.com/razorpay/razorpay-cli/issues" >&2
  exit 1
fi

echo "Checksum verified (${VERSION})."

echo "Extracting..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

# ---------------------------------------------------------------------------
# Install
# ---------------------------------------------------------------------------
mkdir -p "$INSTALL_DIR"
mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

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
