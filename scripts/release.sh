#!/usr/bin/env bash
set -e

cd "$(git rev-parse --show-toplevel)"

# ── pre-flight checks ────────────────────────────────────────────────────────

if ! git diff-index --quiet HEAD --; then
    echo "Error: working tree has uncommitted changes. Commit or stash them first."
    exit 1
fi

current_branch=$(git rev-parse --abbrev-ref HEAD)
git pull origin "${current_branch}"

current_version=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# Suggest the next version by bumping the patch component of the current
# tag. If there is no tag yet, suggest v0.1.0 as the starting point.
suggested_version=""
if [[ "${current_version}" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    major="${BASH_REMATCH[1]}"
    minor="${BASH_REMATCH[2]}"
    patch="${BASH_REMATCH[3]}"
    suggested_version="v${major}.${minor}.$((patch + 1))"
else
    suggested_version="v0.1.0"
fi

echo ""
echo "Current version: ${current_version:-(none)}"
echo "Suggested next:  ${suggested_version}"
echo ""

# ── version ──────────────────────────────────────────────────────────────────

read -rp "Enter new version [${suggested_version}]: " version
version="${version:-${suggested_version}}"

if [[ ! "${version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: invalid version '${version}'. Expected format: vN.N.N"
    exit 1
fi

if git rev-parse "${version}" >/dev/null 2>&1; then
    echo "Error: tag '${version}' already exists."
    exit 1
fi

# ── auto-generated changelog from git log ───────────────────────────────────

tmpchangelog=$(mktemp)
trap 'rm -f "${tmpchangelog}"' EXIT

# Use commit subjects between the previous tag and HEAD as the changelog
# body. Merge commits and the previous release-commit are noise and are
# filtered out. If there is no previous tag, fall back to the full history.
if [ -n "${current_version}" ]; then
    range="${current_version}..HEAD"
else
    range="HEAD"
fi

notes=$(git log "${range}" --no-merges --pretty=format:'- %s' \
        | grep -vE '^- chore: release ' || true)

if [ -z "${notes}" ]; then
    echo "Error: no new commits since ${current_version:-the initial commit}. Nothing to release."
    exit 1
fi

echo ""
echo "── auto-generated notes for ${version} ─────────────────────────────────"
echo "${notes}"
echo "────────────────────────────────────────────────────────────────────────"

# ── write CHANGELOG.md ───────────────────────────────────────────────────────

date=$(date +%Y-%m-%d)

if [ ! -f CHANGELOG.md ]; then
    {
        printf "# Changelog\n\n"
        printf "## %s — %s\n\n%s\n\n" "${version}" "${date}" "${notes}"
    } > "${tmpchangelog}"
else
    {
        printf "# Changelog\n\n"
        printf "## %s — %s\n\n%s\n\n" "${version}" "${date}" "${notes}"
        tail -n +3 CHANGELOG.md
    } > "${tmpchangelog}"
fi

mv "${tmpchangelog}" CHANGELOG.md

# ── update version references in README.md ──────────────────────────────────

# README shows an example of `razorpay --version` output. Keep it in sync
# with the tag we're about to cut so docs do not drift behind releases.
readme_changed=0
if [ -f README.md ] && grep -qE '^razorpay version v[0-9]+\.[0-9]+\.[0-9]+' README.md; then
    # sed -i differs between BSD (macOS) and GNU — handle both.
    if sed --version >/dev/null 2>&1; then
        sed -i -E "s/^razorpay version v[0-9]+\.[0-9]+\.[0-9]+$/razorpay version ${version}/" README.md
    else
        sed -i '' -E "s/^razorpay version v[0-9]+\.[0-9]+\.[0-9]+$/razorpay version ${version}/" README.md
    fi
    if ! git diff --quiet README.md; then
        readme_changed=1
    fi
fi

# ── preview + confirm ────────────────────────────────────────────────────────

echo ""
echo "── CHANGELOG.md preview ────────────────────────────────────────────────"
head -20 CHANGELOG.md
echo "────────────────────────────────────────────────────────────────────────"
if [ "${readme_changed}" -eq 1 ]; then
    echo ""
    echo "README.md version snippet bumped to ${version}."
fi
echo ""

read -rp "Confirm release ${version}? (y/N): " confirm
if [[ ! "${confirm}" =~ ^[Yy]$ ]]; then
    echo "Aborted. Local files have been updated but nothing was committed or tagged."
    echo "Run 'git checkout CHANGELOG.md README.md' to discard the local changes."
    exit 1
fi

# ── commit, tag, push ────────────────────────────────────────────────────────

git add CHANGELOG.md
if [ "${readme_changed}" -eq 1 ]; then
    git add README.md
fi
git commit -m "Release ${version}"
git tag "${version}"
git push origin "${current_branch}"
git push origin "${version}"

echo ""
echo "Released ${version} successfully."
echo "The release.yml workflow will now build and upload the binaries."
