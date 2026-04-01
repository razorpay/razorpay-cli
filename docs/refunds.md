# Refunds

A refund returns captured funds to a customer. Refunds can be full or partial. Partial refunds can be issued multiple times as long as the total does not exceed the captured amount.

## Create a refund

Refunds are created against a specific payment, not directly against an order.

```bash
razorpay refunds create <payment_id> [--amount <amount>] [--speed <speed>] [--param <key=value>]
```

Flags:

| Flag       | Type   | Description                                                        |
| ---------- | ------ | ------------------------------------------------------------------ |
| `--amount` | int    | Amount to refund in smallest currency unit; omit for full refund   |
| `--speed`  | string | Refund processing speed: `normal` or `optimum`                     |
| `--param`  | string | Additional parameter as `key=value`, repeatable                    |

Examples:

```bash
# Full refund
razorpay refunds create pay_29QQoUBi66xm2f

# Partial refund of Rs. 5.00 (500 paise)
razorpay refunds create pay_29QQoUBi66xm2f --amount 500

# Optimum-speed refund (processed within business hours)
razorpay refunds create pay_29QQoUBi66xm2f --speed optimum
```

> `optimum` speed routes refunds through IMPS when possible, resulting in faster settlements compared to `normal`.

## List refunds

```bash
razorpay refunds list [flags]
```

Flags:

| Flag      | Type  | Default | Description                              |
| --------- | ----- | ------- | ---------------------------------------- |
| `--count` | int   | 10      | Number of refunds to fetch (max 100)     |
| `--skip`  | int   | 0       | Records to skip for pagination           |
| `--from`  | int64 | --      | Unix timestamp lower bound (created_at)  |
| `--to`    | int64 | --      | Unix timestamp upper bound (created_at)  |

Example:

```bash
razorpay refunds list --count 50 --from 1704067200
```

## Fetch a refund

```bash
razorpay refunds fetch <refund_id>
```

Example:

```bash
razorpay refunds fetch rfnd_FP8QHiV938haTz
```

## Update a refund

Only the `notes` field can be updated after a refund is created.

```bash
razorpay refunds update <refund_id> --param <key=value> ...
```

Example:

```bash
razorpay refunds update rfnd_FP8QHiV938haTz \
  --param "notes[reason]=Customer requested cancellation"
```

## Refund statuses

| Status     | Meaning                                              |
| ---------- | ---------------------------------------------------- |
| `pending`  | Refund is queued but not yet processed               |
| `processed`| Funds have been returned to the customer             |
| `failed`   | Refund processing failed; funds not returned         |
