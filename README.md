# Razorpay CLI

Command-line interface for the [Razorpay API](https://razorpay.com/docs/api/). Manage payments, orders, customers, invoices, refunds, settlements, disputes, payment links, QR codes, subscriptions, Route, and Smart Collect from your terminal.

## Installation

Recommended — install the latest release for your platform:

```bash
curl -fsSL https://razorpay.com/cli/latest/install.sh | bash
```

Alternatives (`go install`, manual download, build from source) are in [docs/install.md](docs/install.md).

## Configuration

Run `configure` interactively. Existing values are shown in brackets; press Enter to keep them. The secret is masked.

```bash
razorpay configure
```

Or pass credentials via flags — any flag you omit is prompted for:

```bash
razorpay configure --key-id rzp_test_xxxxxxxxxxxx --key-secret xxxxxxxxxxxxxxxxxxxx
```

Credentials are stored in `~/.razorpay/config.yaml`. Environment variables override the file:

```bash
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```

Generate keys from the [Razorpay Dashboard](https://dashboard.razorpay.com/app/website-app-settings/api-keys) — `rzp_test_` for development, `rzp_live_` for production.

## Quick start

Amounts are always in the smallest currency unit (paise for INR).

```bash
# Create an order
razorpay orders create --amount 5000 --currency INR --receipt order-001

# Fetch it back (use the id from the response above)
razorpay orders fetch order_RB58MiP5SPFYyM

# List recent payments
razorpay payments list --count 5
```

## Commands

| Group           | What it manages                                       |
| --------------- | ----------------------------------------------------- |
| `configure`     | Save API credentials                                  |
| `payments`      | Payments — capture, card details, downtime, transfers |
| `orders`        | Orders                                                |
| `customers`     | Customers                                             |
| `refunds`       | Refunds                                               |
| `invoices`      | Invoices and line items                               |
| `payment-links` | Payment Links                                         |
| `qr-codes`      | QR Codes                                              |
| `subscriptions` | Subscriptions and plans                               |
| `route`         | Route — linked accounts and transfers                 |
| `smart-collect` | Smart Collect — virtual accounts                      |
| `settlements`   | Settlements and reconciliation                        |
| `disputes`      | Disputes                                              |
| `documents`     | Documents                                             |

Run `razorpay <group> --help` to list subcommands, or `razorpay <group> <subcommand> --help` for flags and examples.

## Repository layout

```
cmd/
  <resource>/        one package per resource (payments, orders, ...)
    <resource>.go      parent Cmd + subcommand registration
    <subcommand>.go    one file per subcommand (list.go, fetch.go, ...)
  cmdutil/           shared client + error helpers
  configure.go       `configure` command (lives in package cmd directly)
  root.go            root command + subpackage wiring
api/                 HTTP client — auth, JSON pretty-print, multipart upload
config/              config file + env-var loader (viper)
tests/               end-to-end test suite (build tag `e2e`)
docs/                per-resource usage guides
install.sh           platform installer
AGENTS.md            agent guidance for this repo
CHANGELOG.md         release notes
```

To add a subcommand to an existing resource, drop a new file into `cmd/<resource>/` and register it on `Cmd` inside that package's `<resource>.go`. To add a new top-level resource, mirror an existing folder and add `rootCmd.AddCommand(<pkg>.Cmd)` in `cmd/root.go`.

## Tests

The end-to-end suite under [`tests/`](tests/) runs the compiled CLI as a subprocess against the Razorpay **test** API. It is gated by the `e2e` build tag so a plain `go test ./...` skips it.

```bash
export RAZORPAY_TEST_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_TEST_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
go test -tags=e2e -v ./tests/...
```

See [tests/README.md](tests/README.md) for the env variables that unlock destructive subtests (`payments capture`, `refunds create`, `disputes accept`, etc.).

## Documentation

- [docs/install.md](docs/install.md) — platform-specific install
- Per-resource guides: [payments](docs/payments.md) · [orders](docs/orders.md) · [customers](docs/customers.md) · [refunds](docs/refunds.md) · [settlements](docs/settlements.md) · [disputes](docs/disputes.md)
- [CHANGELOG.md](CHANGELOG.md) — release notes
- [AGENTS.md](AGENTS.md) — agent guidance for working in this repository
