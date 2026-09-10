#!/usr/bin/env bash
# update-package-managers.sh — Generate Homebrew formula + Scoop manifest and
# open pull requests for review after a release.
#
# Usage:
#   ./scripts/update-package-managers.sh           # auto-detects latest version from S3
#   ./scripts/update-package-managers.sh 1.0.8     # explicit version override
#
# Prerequisites:
#   - curl
#   - git with push access to the target repos via SSH
#   - gh (GitHub CLI) — authenticated with `gh auth login`

set -euo pipefail

# ---------------------------------------------------------------------------
# Config
# ---------------------------------------------------------------------------
BASE_URL="https://razorpay.com/cli"

# Artifacts are pinned by SHA-256 in both manifests, so they must come from an
# immutable URL. BASE_URL/latest/ is overwritten on every release, which would
# invalidate the pinned hash the moment the next version ships; GitHub release
# assets are per-tag and never change.
RELEASE_BASE="https://github.com/razorpay/razorpay-cli/releases/download"

HOMEBREW_REPO="git@github.com:razorpay/homebrew-razorpay-cli.git"
SCOOP_REPO="git@github.com:razorpay/scoop-razorpay-cli.git"

# ---------------------------------------------------------------------------
# Resolve version — auto-detect from S3 or accept as argument
# ---------------------------------------------------------------------------
if [[ $# -ge 1 ]]; then
  VERSION="$1"
else
  echo "==> Auto-detecting latest version from ${BASE_URL}/latest/version..."
  VERSION=$(curl -fsSL "${BASE_URL}/latest/version")
fi

VERSION_NUM="${VERSION#v}"  # strip leading v if present

if [[ -z "${VERSION_NUM}" ]]; then
  echo "ERROR: Could not determine version. Pass it manually: $0 <version>"
  exit 1
fi

echo "==> Updating package managers for version ${VERSION_NUM}"

# ---------------------------------------------------------------------------
# Download checksums
# ---------------------------------------------------------------------------
TMPDIR=$(mktemp -d)
trap 'rm -rf "${TMPDIR}"' EXIT

# Read the checksums for the tag being published rather than from latest/. The
# version is resolved from latest/version but the hashes must describe *that*
# version, and the two prefixes can drift apart during a partial promotion.
echo "==> Downloading checksums for v${VERSION_NUM}..."
curl -fsSL "${RELEASE_BASE}/v${VERSION_NUM}/razorpay-mac-checksums.txt" -o "${TMPDIR}/mac-checksums.txt"
curl -fsSL "${RELEASE_BASE}/v${VERSION_NUM}/razorpay-windows-checksums.txt" -o "${TMPDIR}/windows-checksums.txt"

# ---------------------------------------------------------------------------
# Extract SHA256 hashes
# ---------------------------------------------------------------------------
# Match the whole filename. A substring match would silently return two hashes
# on separate lines if a future archive name ever contained another as a suffix.
get_sha() {
  local file="$1" checksums="$2" sha
  sha=$(awk -v name="${file}" '$2 == name { print $1 }' "${checksums}")
  if [[ -z "${sha}" ]]; then
    echo "ERROR: ${checksums##*/} has no entry for ${file}" >&2
    exit 1
  fi
  echo "${sha}"
}

MAC_ARM64_FILE="razorpay_${VERSION_NUM}_mac-os_arm64.tar.gz"
MAC_AMD64_FILE="razorpay_${VERSION_NUM}_mac-os_x86_64.tar.gz"
WIN_AMD64_FILE="razorpay_${VERSION_NUM}_windows_x86_64.zip"
WIN_I386_FILE="razorpay_${VERSION_NUM}_windows_i386.zip"

MAC_ARM64_SHA=$(get_sha "${MAC_ARM64_FILE}" "${TMPDIR}/mac-checksums.txt")
MAC_AMD64_SHA=$(get_sha "${MAC_AMD64_FILE}" "${TMPDIR}/mac-checksums.txt")
WIN_AMD64_SHA=$(get_sha "${WIN_AMD64_FILE}" "${TMPDIR}/windows-checksums.txt")
WIN_I386_SHA=$(get_sha "${WIN_I386_FILE}" "${TMPDIR}/windows-checksums.txt")

echo "  mac arm64:     ${MAC_ARM64_SHA}"
echo "  mac x86_64:    ${MAC_AMD64_SHA}"
echo "  win x86_64:    ${WIN_AMD64_SHA}"
echo "  win i386:      ${WIN_I386_SHA}"

# ---------------------------------------------------------------------------
# Verify every URL that is about to be pinned
#
# The manifests pin a hash to a URL. If the two ever disagree, `brew install`
# and `scoop install` fail outright for everyone, so confirm here -- before a
# PR is opened -- that each URL resolves and serves exactly the bytes whose
# hash is going into the manifest.
# ---------------------------------------------------------------------------
sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{ print $1 }'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{ print $1 }'
  else
    echo "ERROR: neither sha256sum nor shasum is available." >&2
    exit 1
  fi
}

verify_pinned_url() {
  local url="$1" expected="$2" out actual
  out="${TMPDIR}/$(basename "${url}")"

  if ! curl -fsSL "${url}" -o "${out}"; then
    echo "ERROR: ${url} did not resolve." >&2
    echo "       Release assets are published by .goreleaser/*.yml; check the tag was built." >&2
    exit 1
  fi

  actual=$(sha256_of "${out}")
  if [[ "${expected}" != "${actual}" ]]; then
    echo "ERROR: ${url}" >&2
    echo "       expected ${expected}" >&2
    echo "       actual   ${actual}" >&2
    echo "       Refusing to generate a manifest that cannot install." >&2
    exit 1
  fi
  echo "  ok  $(basename "${url}")"
}

echo "==> Verifying pinned URLs..."
verify_pinned_url "${RELEASE_BASE}/v${VERSION_NUM}/${MAC_ARM64_FILE}" "${MAC_ARM64_SHA}"
verify_pinned_url "${RELEASE_BASE}/v${VERSION_NUM}/${MAC_AMD64_FILE}" "${MAC_AMD64_SHA}"
verify_pinned_url "${RELEASE_BASE}/v${VERSION_NUM}/${WIN_AMD64_FILE}" "${WIN_AMD64_SHA}"
verify_pinned_url "${RELEASE_BASE}/v${VERSION_NUM}/${WIN_I386_FILE}" "${WIN_I386_SHA}"

# ---------------------------------------------------------------------------
# Generate Homebrew formula
# ---------------------------------------------------------------------------
echo "==> Generating Homebrew formula..."

FORMULA_DIR="${TMPDIR}/homebrew"
mkdir -p "${FORMULA_DIR}"

cat > "${FORMULA_DIR}/razorpay.rb" <<RUBY
# typed: false
# frozen_string_literal: true

# This file was generated by scripts/update-package-managers.sh. Do not edit manually.
class Razorpay < Formula
  desc "Official Razorpay CLI."
  homepage "https://github.com/razorpay/razorpay-cli"
  version "${VERSION_NUM}"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "${RELEASE_BASE}/v${VERSION_NUM}/${MAC_ARM64_FILE}"
      sha256 "${MAC_ARM64_SHA}"
    else
      url "${RELEASE_BASE}/v${VERSION_NUM}/${MAC_AMD64_FILE}"
      sha256 "${MAC_AMD64_SHA}"
    end
  end

  def install
    bin.install "razorpay"
  end

  def caveats
    <<~EOS
      Thanks for installing the Razorpay CLI! If this is your first time using the CLI, be sure to run \`razorpay configure\` first.
    EOS
  end

  test do
    assert_match "razorpay version", shell_output("#{bin}/razorpay --version")
  end
end
RUBY

echo "  Formula written to ${FORMULA_DIR}/razorpay.rb"

# ---------------------------------------------------------------------------
# Generate Scoop manifest
# ---------------------------------------------------------------------------
echo "==> Generating Scoop manifest..."

SCOOP_DIR="${TMPDIR}/scoop"
mkdir -p "${SCOOP_DIR}"

cat > "${SCOOP_DIR}/razorpay.json" <<JSON
{
  "version": "${VERSION_NUM}",
  "description": "Official Razorpay CLI.",
  "homepage": "https://github.com/razorpay/razorpay-cli",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "${RELEASE_BASE}/v${VERSION_NUM}/${WIN_AMD64_FILE}",
      "hash": "${WIN_AMD64_SHA}"
    },
    "32bit": {
      "url": "${RELEASE_BASE}/v${VERSION_NUM}/${WIN_I386_FILE}",
      "hash": "${WIN_I386_SHA}"
    }
  },
  "bin": "razorpay.exe",
  "checkver": {
    "url": "${BASE_URL}/latest/version",
    "regex": "v([\\\\d.]+)"
  },
  "autoupdate": {
    "architecture": {
      "64bit": {
        "url": "${RELEASE_BASE}/v\$version/razorpay_\$version_windows_x86_64.zip"
      },
      "32bit": {
        "url": "${RELEASE_BASE}/v\$version/razorpay_\$version_windows_i386.zip"
      }
    },
    "hash": {
      "url": "${RELEASE_BASE}/v\$version/razorpay-windows-checksums.txt"
    }
  }
}
JSON

echo "  Manifest written to ${SCOOP_DIR}/razorpay.json"

# ---------------------------------------------------------------------------
# Create PR for Homebrew formula
# ---------------------------------------------------------------------------
echo "==> Creating PR for Homebrew formula in ${HOMEBREW_REPO}..."

BREW_CLONE="${TMPDIR}/homebrew-razorpay-cli"
git clone "${HOMEBREW_REPO}" "${BREW_CLONE}" 2>/dev/null || {
  echo "ERROR: Could not clone ${HOMEBREW_REPO}. Make sure the repo exists and you have access."
  exit 1
}

BRANCH_NAME="update/v${VERSION_NUM}"
HOMEBREW_PR_URL=""
SCOOP_PR_URL=""

cd "${BREW_CLONE}"
# Switch to the branch — reuse it if it already exists on the remote
if git ls-remote --exit-code --heads origin "${BRANCH_NAME}" >/dev/null 2>&1; then
  git fetch origin "${BRANCH_NAME}"
  git checkout "${BRANCH_NAME}"
else
  git checkout -b "${BRANCH_NAME}"
fi

# Copy the generated formula after switching branches to avoid checkout conflicts
cp "${FORMULA_DIR}/razorpay.rb" "${BREW_CLONE}/razorpay.rb"
git add razorpay.rb
if git diff --cached --quiet; then
  echo "  No changes to Homebrew formula — skipping."
else
  git commit -m "Update razorpay formula to v${VERSION_NUM}"
  git push -u origin "${BRANCH_NAME}"
  # Open a PR (skip if an open one already exists for this branch)
  if HOMEBREW_PR_URL=$(gh pr view "${BRANCH_NAME}" --repo razorpay/homebrew-razorpay-cli --json url,state --jq 'select(.state == "OPEN") | .url' 2>/dev/null) && [[ -n "${HOMEBREW_PR_URL}" ]]; then
    echo "  PR already exists — updated branch pushed."
  else
    HOMEBREW_PR_URL=$(gh pr create \
      --repo razorpay/homebrew-razorpay-cli \
      --head "${BRANCH_NAME}" \
      --title "Update razorpay formula to v${VERSION_NUM}" \
      --body "Auto-generated by \`scripts/update-package-managers.sh\`. Updates the Homebrew formula to v${VERSION_NUM}.")
  fi
  echo "  Homebrew PR: ${HOMEBREW_PR_URL}"
fi

# ---------------------------------------------------------------------------
# Create PR for Scoop manifest
# ---------------------------------------------------------------------------
echo "==> Creating PR for Scoop manifest in ${SCOOP_REPO}..."

SCOOP_CLONE="${TMPDIR}/scoop-razorpay-cli"
git clone "${SCOOP_REPO}" "${SCOOP_CLONE}" 2>/dev/null || {
  echo "ERROR: Could not clone ${SCOOP_REPO}. Make sure the repo exists and you have access."
  exit 1
}

cd "${SCOOP_CLONE}"
# Switch to the branch — reuse it if it already exists on the remote
if git ls-remote --exit-code --heads origin "${BRANCH_NAME}" >/dev/null 2>&1; then
  git fetch origin "${BRANCH_NAME}"
  git checkout "${BRANCH_NAME}"
else
  git checkout -b "${BRANCH_NAME}"
fi

# Copy the generated manifest after switching branches to avoid checkout conflicts
cp "${SCOOP_DIR}/razorpay.json" "${SCOOP_CLONE}/razorpay.json"
git add razorpay.json
if git diff --cached --quiet; then
  echo "  No changes to Scoop manifest — skipping."
else
  git commit -m "Update razorpay manifest to v${VERSION_NUM}"
  git push -u origin "${BRANCH_NAME}"
  # Open a PR (skip if an open one already exists for this branch)
  if SCOOP_PR_URL=$(gh pr view "${BRANCH_NAME}" --repo razorpay/scoop-razorpay-cli --json url,state --jq 'select(.state == "OPEN") | .url' 2>/dev/null) && [[ -n "${SCOOP_PR_URL}" ]]; then
    echo "  PR already exists — updated branch pushed."
  else
    SCOOP_PR_URL=$(gh pr create \
      --repo razorpay/scoop-razorpay-cli \
      --head "${BRANCH_NAME}" \
      --title "Update razorpay manifest to v${VERSION_NUM}" \
      --body "Auto-generated by \`scripts/update-package-managers.sh\`. Updates the Scoop manifest to v${VERSION_NUM}.")
  fi
  echo "  Scoop PR: ${SCOOP_PR_URL}"
fi

# ---------------------------------------------------------------------------
echo ""
echo "Done! Review and merge the PRs to publish v${VERSION_NUM}."
echo ""
if [[ -n "${HOMEBREW_PR_URL}" ]]; then
  echo "  Homebrew PR: ${HOMEBREW_PR_URL}"
fi
if [[ -n "${SCOOP_PR_URL}" ]]; then
  echo "  Scoop PR:    ${SCOOP_PR_URL}"
fi
echo ""
echo "Once merged, users can install via:"
echo "  brew install razorpay/razorpay-cli/razorpay    (macOS)"
echo "  scoop bucket add razorpay https://github.com/razorpay/scoop-razorpay-cli && scoop install razorpay   (Windows)"
