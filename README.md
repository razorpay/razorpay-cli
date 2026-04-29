# Razorpay CLI

A command-line interface for the [Razorpay API](https://razorpay.com/docs/api/). Interact with payments, orders, customers, refunds, settlements, and disputes directly from your terminal.

## Installation

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.sh | sh
```

The script detects your OS and architecture, downloads the appropriate binary from the [latest GitHub release](https://github.com/razorpay/razorpay-cli/releases/latest), verifies the SHA-256 checksum, and installs the binary to `~/.local/bin` (or `$XDG_BIN_HOME` if set). Running as root installs to `/usr/local/bin`.

To override the install directory:

```bash
RAZORPAY_INSTALL=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.sh | sh
```

### Windows

Run the following in PowerShell:

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.ps1 | iex"
```

The script downloads the latest release, verifies the checksum, installs to `%USERPROFILE%\.local\bin`, and adds that directory to your user `PATH`.

To override the install directory:

```powershell
$env:RAZORPAY_INSTALL = "C:\tools\bin"
powershell -ExecutionPolicy ByPass -c "irm https://raw.githubusercontent.com/razorpay/razorpay-cli/master/install.ps1 | iex"
```

### Manual download

Pre-built binaries for all platforms are available on the [releases page](https://github.com/razorpay/razorpay-cli/releases). Download the archive for your platform, extract the `razorpay` binary, and place it somewhere on your `PATH`.

### Install with `go install`

```bash
go install github.com/razorpay/razorpay-cli@latest
```

Requires Go 1.21 or later.

### Install from source

```bash
git clone https://github.com/razorpay/razorpay-cli.git
cd razorpay-cli
go build -o razorpay .
mv razorpay /usr/local/bin/razorpay
```

## Configuration

### Interactive setup

Run the configure command and enter your API Key ID and Key Secret when prompted. The key secret input is masked.

```bash
razorpay configure
```

Credentials are saved to `~/.razorpay/config.yaml`.

### Environment variables

Set credentials via environment variables to override the config file. This is useful in CI/CD environments.

```bash
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```

You can generate API keys from the [Razorpay Dashboard](https://dashboard.razorpay.com/app/website-app-settings/api-keys). Use `rzp_test_` keys for development and `rzp_live_` keys for production.

## Hello World

The following walkthrough creates an order and then fetches it back. These are the two most fundamental operations in the Razorpay payment flow.

### Step 1 -- Configure credentials

```bash
razorpay configure
# Key ID: rzp_test_xxxxxxxxxxxx
# Key Secret: (hidden)
# Credentials saved to /Users/you/.razorpay/config.yaml
```

### Step 2 -- Create an order

An order represents a payment intent. Amounts are always in the smallest currency unit (paise for INR).

```bash
razorpay orders create --amount 5000 --currency INR --receipt "order-001"
```

Expected output:

```json
{
  "id": "order_RB58MiP5SPFYyM",
  "entity": "order",
  "amount": 5000,
  "amount_paid": 0,
  "amount_due": 5000,
  "currency": "INR",
  "receipt": "order-001",
  "status": "created",
  "attempts": 0,
  "created_at": 1756455561
}
```

### Step 3 -- Fetch the order

Use the `id` from the previous response to fetch the order.

```bash
razorpay orders fetch order_RB58MiP5SPFYyM
```

### Step 4 -- List recent payments

Once a customer completes payment against the order, it appears in the payments list.

```bash
razorpay payments list --count 5
```

## Available Commands

| Command                                 | Description                              |
| --------------------------------------- | ---------------------------------------- |
| `razorpay configure`                    | Save API credentials to config file      |
| `razorpay payments list`                | List payments                            |
| `razorpay payments fetch <id>`          | Fetch a payment by ID                    |
| `razorpay payments capture <id>`        | Capture an authorized payment            |
| `razorpay payments update <id>`         | Update payment metadata                  |
| `razorpay payments transfers <id>`      | Fetch transfers for a payment            |
| `razorpay orders list`                  | List orders                              |
| `razorpay orders fetch <id>`            | Fetch an order by ID                     |
| `razorpay orders create`                | Create a new order                       |
| `razorpay orders update <id>`           | Update an order                          |
| `razorpay orders payments <id>`         | Fetch payments for an order              |
| `razorpay customers list`               | List customers                           |
| `razorpay customers fetch <id>`         | Fetch a customer by ID                   |
| `razorpay customers create`             | Create a new customer                    |
| `razorpay customers update <id>`        | Update a customer                        |
| `razorpay refunds list`                 | List refunds                             |
| `razorpay refunds fetch <id>`           | Fetch a refund by ID                     |
| `razorpay refunds create <payment_id>`  | Create a refund for a payment            |
| `razorpay refunds update <id>`          | Update a refund                          |
| `razorpay settlements list`             | List settlements                         |
| `razorpay settlements fetch <id>`       | Fetch a settlement by ID                 |
| `razorpay settlements recon`            | Fetch settlement recon report            |
| `razorpay disputes list`                | List disputes                            |
| `razorpay disputes fetch <id>`          | Fetch a dispute by ID                    |
| `razorpay disputes accept <id>`         | Accept a dispute                         |
| `razorpay disputes contest <id>`        | Contest a dispute                        |

Run `razorpay [command] --help` for flags and usage details on any command.

## Documentation

Detailed usage guides for each resource are in the [docs/](docs/) directory.

- [Payments](docs/payments.md)
- [Orders](docs/orders.md)
- [Customers](docs/customers.md)
- [Refunds](docs/refunds.md)
- [Settlements](docs/settlements.md)
- [Disputes](docs/disputes.md)
