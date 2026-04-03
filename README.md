# Razorpay CLI

A command-line interface for the [Razorpay API](https://razorpay.com/docs/api/). Interact with payments, orders, customers, refunds, settlements, and disputes directly from your terminal.

## Installation

```bash
go install github.com/razorpay/razorpay-cli@latest
```

### Prerequisites

- Go 1.21 or later

### Install from source

```bash
git clone https://github.com/razorpay/razorpay-cli.git
cd razorpay-cli
go build -o razorpay .
```

Move the binary somewhere on your `PATH`:

```bash
mv razorpay /usr/local/bin/razorpay
```

### Install with `go install`

```bash
go install github.com/razorpay/razorpay-cli@latest
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
