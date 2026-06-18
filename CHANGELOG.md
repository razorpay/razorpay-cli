# Changelog

## v1.0.9 — 2026-06-18

- feat: nudge users when a new release is available
- feat: identify CLI to API via User-Agent header

## v1.0.8 — 2026-05-22

- Added missing query params and removed undocumented ones across orders, disputes, settlements, and subscriptions - Replaced hardcoded values with configurable flags in route and smart-collect commands 
- Accepted JSON input for complex nested arrays (transfers, allowed_payers, others) 
- Removed transfer_fetch_by_settlement and merged its functionality into transfer list 
- Ensured consistent flag naming between create and update commands

## v1.0.7 — 2026-05-11

- Removed sudo command from the install.sh script

## v1.0.6 — 2026-05-11

- install.sh: Fixed broken S3_BASE_URL that caused 404s; simplified to download directly from latest/ folder 
- docs/install.md: Rewrote to match production docs structure with curl commands, quick-install one-liner, macOS quarantine step, and credential configuration

## v1.0.5 — 2026-05-11

- Enhance make release flow with auto-generated notes and version bump
- Drop yml alias and keep only yaml as the format name
- Surface make targets in README docs instead of raw go commands
- Unify configure prompts so output format reads from stdin uniformly
- Add configurable output format with JSON, YAML, and TOML
- Fix e2e workflow tests against capability-limited accounts
- Wire make test and ground README in real examples
- Restructure e2e tests as chained lifecycle workflows
- Refresh README for post-merge command surface
- Polish CLI UX and add end-to-end test suite

## v1.0.4 — 2026-04-30

Added:
- Non-interactive razorpay configure --key-id <id> --key-secret <secret>
Fixed:
-Incorrect example in razorpay help output

## v1.0.3 — 2026-04-29

-Checking the cloudfront changes and ensuring the correct version for razorpay cli.

## v1.0.2 — 2026-04-29

-Added changes to the release.yml file.

## v1.0.1 — 2026-04-29

- Added missing version number for the integrated apis.

## v1.0.0 — 2026-04-28

- Initial release of the Razorpay CLI
- Commands: payments, orders, refunds, customers, invoices, subscriptions, payment links, disputes, settlements, routes, smart collect, QR codes, documents
- Configure API credentials via `razorpay configure` or environment variables
- JSON output for all API responses

