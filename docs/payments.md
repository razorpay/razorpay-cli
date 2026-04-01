# Payments

Payments represent a customer's payment attempt against an order or a direct charge. The Razorpay payment lifecycle is: `created` -> `authorized` -> `captured` -> `refunded`.

## List payments

Fetches a paginated list of payments in reverse chronological order.

```bash
razorpay payments list [flags]
```

Flags:

| Flag        | Type  | Default | Description                              |
| ----------- | ----- | ------- | ---------------------------------------- |
| `--count`   | int   | 10      | Number of payments to fetch (max 100)    |
| `--skip`    | int   | 0       | Number of records to skip for pagination |
| `--from`    | int64 | --      | Unix timestamp lower bound (created_at)  |
| `--to`      | int64 | --      | Unix timestamp upper bound (created_at)  |

Examples:

```bash
# Fetch the 20 most recent payments
razorpay payments list --count 20

# Fetch payments from January 2024
razorpay payments list --from 1704067200 --to 1706745600
```

## Fetch a payment

```bash
razorpay payments fetch <payment_id>
```

Example:

```bash
razorpay payments fetch pay_29QQoUBi66xm2f
```

## Capture a payment

Changes a payment's status from `authorized` to `captured`. The `--amount` must equal the authorized amount.

```bash
razorpay payments capture <payment_id> --amount <amount> [--currency <code>]
```

Flags:

| Flag         | Type   | Default | Description                                       |
| ------------ | ------ | ------- | ------------------------------------------------- |
| `--amount`   | int    | --      | Amount in smallest currency unit, e.g. paise      |
| `--currency` | string | `INR`   | ISO 4217 currency code                            |

Example:

```bash
# Capture INR 10.00 (1000 paise)
razorpay payments capture pay_29QQoUBi66xm2f --amount 1000 --currency INR
```

> Attempting to capture a payment that is not in the `authorized` state returns a `400` error.

## Update a payment

Updates the `notes` field on a payment. Notes are key-value pairs stored for your internal reference.

```bash
razorpay payments update <payment_id> --param <key=value> ...
```

Example:

```bash
razorpay payments update pay_29QQoUBi66xm2f \
  --param "notes[order_ref]=ORD-1234" \
  --param "notes[customer_ref]=CUST-5678"
```

## Fetch transfers for a payment

Returns all Route transfers created from a payment.

```bash
razorpay payments transfers <payment_id>
```

Example:

```bash
razorpay payments transfers pay_29QQoUBi66xm2f
```

## Payment statuses

| Status       | Meaning                                              |
| ------------ | ---------------------------------------------------- |
| `created`    | Payment initiated but not yet authorized             |
| `authorized` | Bank has authorized; capture required                |
| `captured`   | Amount deducted from customer; funds with Razorpay   |
| `refunded`   | Amount refunded partially or fully                   |
| `failed`     | Payment failed at authorization or capture           |
