#!/bin/sh
# install.sh - Installer for the Razorpay CLI
# Usage: curl -fsSL https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.sh | sh
set -eu

BINARY_NAME="razorpay"
REPO="razorpay/razorpay-cli"
GITHUB_API="https://api.github.com/repos/${REPO}"
GITHUB_RELEASES="${GITHUB_API}/releases/latest"

# ---- helpers ----------------------------------------------------------------

say() {
    printf 'razorpay-cli: %s\n' "$1"
}

err() {
    say "error: $1" >&2
    exit 1
}

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        err "required command not found: '$1'"
    fi
}

# ---- platform detection -----------------------------------------------------

detect_os() {
    OS="$(uname -s)"
    case "${OS}" in
        Linux)  echo "Linux" ;;
        Darwin) echo "Darwin" ;;
        *)      err "unsupported operating system: ${OS}" ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "${ARCH}" in
        x86_64 | amd64)          echo "x86_64" ;;
        aarch64 | arm64)         echo "arm64" ;;
        i386 | i686)             echo "i386" ;;
        *)                       err "unsupported architecture: ${ARCH}" ;;
    esac
}

# ---- install directory selection --------------------------------------------

# Follows XDG Base Directory spec then falls back to ~/.local/bin.
# /usr/local/bin is used only when running as root.
resolve_install_dir() {
    if [ -n "${RAZORPAY_INSTALL:-}" ]; then
        echo "${RAZORPAY_INSTALL}"
    elif [ -n "${XDG_BIN_HOME:-}" ]; then
        echo "${XDG_BIN_HOME}"
    elif [ "$(id -u)" -eq 0 ]; then
        echo "/usr/local/bin"
    else
        echo "${HOME}/.local/bin"
    fi
}

# ---- download ---------------------------------------------------------------

download() {
    URL="$1"
    DEST="$2"
    if command -v curl >/dev/null 2>&1; then
        curl --proto '=https' --tlsv1.2 -fsSL "${URL}" -o "${DEST}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${DEST}" "${URL}"
    else
        err "neither curl nor wget found; please install one and retry"
    fi
}

# ---- checksum verification --------------------------------------------------

verify_checksum() {
    ARCHIVE="$1"
    CHECKSUMS_FILE="$2"
    ARCHIVE_NAME="$(basename "${ARCHIVE}")"

    if command -v sha256sum >/dev/null 2>&1; then
        EXPECTED="$(grep "${ARCHIVE_NAME}" "${CHECKSUMS_FILE}" | awk '{print $1}')"
        ACTUAL="$(sha256sum "${ARCHIVE}" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
        EXPECTED="$(grep "${ARCHIVE_NAME}" "${CHECKSUMS_FILE}" | awk '{print $1}')"
        ACTUAL="$(shasum -a 256 "${ARCHIVE}" | awk '{print $1}')"
    else
        say "warning: no sha256 tool found, skipping checksum verification"
        return 0
    fi

    if [ "${EXPECTED}" != "${ACTUAL}" ]; then
        err "checksum mismatch for ${ARCHIVE_NAME}\n  expected: ${EXPECTED}\n  actual:   ${ACTUAL}"
    fi
    say "checksum verified"
}

# ---- main -------------------------------------------------------------------

main() {
    need_cmd uname

    OS="$(detect_os)"
    ARCH="$(detect_arch)"
    INSTALL_DIR="$(resolve_install_dir)"

    say "detecting latest release..."
    if command -v curl >/dev/null 2>&1; then
        LATEST_TAG="$(curl --proto '=https' --tlsv1.2 -fsSL "${GITHUB_RELEASES}" \
            | grep '"tag_name"' \
            | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
    elif command -v wget >/dev/null 2>&1; then
        LATEST_TAG="$(wget -qO- "${GITHUB_RELEASES}" \
            | grep '"tag_name"' \
            | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
    else
        err "neither curl nor wget found; please install one and retry"
    fi

    if [ -z "${LATEST_TAG}" ]; then
        err "could not determine latest release tag"
    fi

    # Strip leading 'v' for the version number used in filenames
    VERSION="${LATEST_TAG#v}"

    ARCHIVE_NAME="razorpay-cli_${OS}_${ARCH}.tar.gz"
    CHECKSUMS_NAME="razorpay-cli_${VERSION}_checksums.txt"
    BASE_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}"

    say "installing razorpay-cli ${LATEST_TAG} (${OS}/${ARCH})"

    TMP_DIR="$(mktemp -d)"
    trap 'rm -rf "${TMP_DIR}"' EXIT

    say "downloading ${ARCHIVE_NAME}..."
    download "${BASE_URL}/${ARCHIVE_NAME}"    "${TMP_DIR}/${ARCHIVE_NAME}"
    download "${BASE_URL}/${CHECKSUMS_NAME}" "${TMP_DIR}/${CHECKSUMS_NAME}"

    verify_checksum "${TMP_DIR}/${ARCHIVE_NAME}" "${TMP_DIR}/${CHECKSUMS_NAME}"

    say "extracting..."
    tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

    # Locate the binary — GoReleaser may or may not wrap in a subdirectory
    BINARY_PATH="$(find "${TMP_DIR}" -type f -name "${BINARY_NAME}" | head -1)"
    if [ -z "${BINARY_PATH}" ]; then
        err "could not find '${BINARY_NAME}' binary in the downloaded archive"
    fi

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    say "installing to ${INSTALL_DIR}/${BINARY_NAME}"
    cp "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    # Verify the binary works
    if ! "${INSTALL_DIR}/${BINARY_NAME}" --help >/dev/null 2>&1; then
        err "installed binary failed to run — please report this at https://github.com/${REPO}/issues"
    fi

    say "razorpay-cli ${LATEST_TAG} installed successfully!"
    say ""

    # PATH guidance
    case ":${PATH}:" in
        *":${INSTALL_DIR}:"*)
            # Already on PATH — nothing to do
            ;;
        *)
            say "add the following to your shell profile to put razorpay on your PATH:"
            say ""
            say "  export PATH=\"${INSTALL_DIR}:\$PATH\""
            say ""
            say "then restart your shell or run:  source ~/.bashrc  (or ~/.zshrc)"
            ;;
    esac
}

main "$@"
