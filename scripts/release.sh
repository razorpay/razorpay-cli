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

last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "(none)")
echo "Last release: ${last_tag}"
echo ""

# ── version ──────────────────────────────────────────────────────────────────

read -rp "Enter new version (format: vN.N.N): " version

if [[ ! "${version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: invalid version '${version}'. Expected format: vN.N.N"
    exit 1
fi

if git rev-parse "${version}" >/dev/null 2>&1; then
    echo "Error: tag '${version}' already exists."
    exit 1
fi

# ── changelog entry (open $EDITOR) ───────────────────────────────────────────

tmpfile=$(mktemp /tmp/razorpay-changelog-XXXXXX.md)
trap 'rm -f "${tmpfile}"' EXIT

cat > "${tmpfile}" <<EOF
# Changelog entry for ${version}
# Lines starting with '#' will be ignored.
# Use bullet points for each entry, for example:
#
#   - feat: Added payment capture command
#   - fix: Fixed auth token expiry on refresh
#
# Save and quit the editor when done. An empty entry will abort the release.

EOF

${EDITOR:-vi} "${tmpfile}"

# Strip comment lines and collapse blank lines
notes=$(grep -v '^[[:space:]]*#' "${tmpfile}" | sed '/^[[:space:]]*$/d')

if [ -z "${notes}" ]; then
    echo "Changelog is empty. Aborting."
    exit 1
fi

# ── write CHANGELOG.md ───────────────────────────────────────────────────────

date=$(date +%Y-%m-%d)
tmpchangelog=$(mktemp)
trap 'rm -f "${tmpfile}" "${tmpchangelog}"' EXIT

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

echo ""
echo "── CHANGELOG.md preview ────────────────────────────────────────────────"
head -20 CHANGELOG.md
echo "────────────────────────────────────────────────────────────────────────"
echo ""

read -rp "Confirm release ${version}? (y/N): " confirm
if [[ ! "${confirm}" =~ ^[Yy]$ ]]; then
    echo "Aborted. CHANGELOG.md has been updated locally but nothing was committed or tagged."
    echo "Run 'git checkout CHANGELOG.md' to discard the local change."
    exit 1
fi

# ── commit, tag, push ────────────────────────────────────────────────────────

git add CHANGELOG.md
git commit -m "chore: release ${version}"
git tag "${version}"
git push origin "${current_branch}"
git push origin "${version}"

echo ""
echo "Released ${version} successfully."
echo "The release.yml workflow will now build and upload the binaries."
