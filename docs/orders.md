# Orders

An order is a payment intent that you create before presenting the Razorpay checkout to your customer. Every payment in Razorpay is recommended to be associated with an order for full payment tracking.

## Create an order

```bash
razorpay orders create --amount <amount> [--currency <code>] [--receipt <id>] [--param <key=value>]
```

Flags:

| Flag          | Type   | Default | Description                                        |
| ------------- | ------ | ------- | -------------------------------------------------- |
| `--amount`    | int    | --      | Amount in smallest currency unit (required)        |
| `--currency`  | string | `INR`   | ISO 4217 currency code                             |
| `--receipt`   | string | --      | Your internal receipt or order ID (max 40 chars)   |
| `--param`     | string | --      | Additional parameter as `key=value`, repeatable    |

Examples:

```bash
# Create an INR order for Rs. 299 (29900 paise)
razorpay orders create --amount 29900 --currency INR --receipt "receipt-001"

# Create an order with partial payment enabled
razorpay orders create --amount 50000 --currency INR --param partial_payment=true
```

## List orders

```bash
razorpay orders list [flags]
```

Flags:

| Flag        | Type   | Default | Description                                     |
| ----------- | ------ | ------- | ----------------------------------------------- |
| `--count`   | int    | 10      | Number of orders to fetch (max 100)              |
| `--skip`    | int    | 0       | Records to skip for pagination                  |
| `--from`    | int64  | --      | Unix timestamp lower bound                      |
| `--to`      | int64  | --      | Unix timestamp upper bound                      |
| `--status`  | string | --      | Filter by status: `created`, `attempted`, `paid` |

Examples:

```bash
# List all unpaid orders
razorpay orders list --status created

# List 50 orders created this week
razorpay orders list --count 50 --from 1711929600
```

## Fetch an order

```bash
razorpay orders fetch <order_id>
```

Example:

```bash
razorpay orders fetch order_RB58MiP5SPFYyM
```

## Update an order

Only the `notes` field can be updated after order creation.

```bash
razorpay orders update <order_id> --param <key=value> ...
```

Example:

```bash
razorpay orders update order_RB58MiP5SPFYyM \
  --param "notes[shipping_id]=SHIP-9012"
```

## Fetch payments for an order

Returns all payment attempts (successful and failed) made against an order.

```bash
razorpay orders payments <order_id>
```

Example:

```bash
razorpay orders payments order_RB58MiP5SPFYyM
```

## Order statuses

| Status      | Meaning                                                        |
| ----------- | -------------------------------------------------------------- |
| `created`   | No payment attempt has been made yet                           |
| `attempted` | At least one payment attempt has been made; none captured yet  |
| `paid`      | A payment has been captured; no further payments are accepted  |
