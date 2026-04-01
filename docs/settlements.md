# Settlements

Settlements represent the transfer of funds from Razorpay to your bank account. Razorpay batches captured payments and settles them on a T+2 or T+3 cycle depending on your plan.

## List settlements

```bash
razorpay settlements list [flags]
```

Flags:

| Flag      | Type  | Default | Description                              |
| --------- | ----- | ------- | ---------------------------------------- |
| `--count` | int   | 10      | Number of settlements to fetch           |
| `--skip`  | int   | 0       | Records to skip for pagination           |
| `--from`  | int64 | --      | Unix timestamp lower bound (created_at)  |
| `--to`    | int64 | --      | Unix timestamp upper bound (created_at)  |

Example:

```bash
# Settlements from Q1 2024
razorpay settlements list --from 1704067200 --to 1711929599 --count 100
```

## Fetch a settlement

```bash
razorpay settlements fetch <settlement_id>
```

Example:

```bash
razorpay settlements fetch setl_7IbFiDRTBaQqeL
```

## Fetch settlement recon report

The recon report contains a line-by-line breakdown of every transaction included in a settlement. Use this for reconciling settlements against your internal records.

```bash
razorpay settlements recon [--year <year>] [--month <month>] [--day <day>]
```

Flags:

| Flag      | Type | Description                          |
| --------- | ---- | ------------------------------------ |
| `--year`  | int  | Year, e.g. `2024`                    |
| `--month` | int  | Month as a number, `1` through `12`  |
| `--day`   | int  | Day of the month, `1` through `31`   |

Examples:

```bash
# Recon report for January 2024
razorpay settlements recon --year 2024 --month 1

# Recon report for a specific day
razorpay settlements recon --year 2024 --month 3 --day 15
```
