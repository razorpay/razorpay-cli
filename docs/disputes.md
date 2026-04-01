# Disputes

A dispute is raised when a customer files a chargeback with their bank. Razorpay notifies you via webhook and gives you a window to either accept the dispute or contest it with evidence.

## List disputes

```bash
razorpay disputes list [flags]
```

Flags:

| Flag      | Type  | Default | Description                             |
| --------- | ----- | ------- | --------------------------------------- |
| `--count` | int   | 10      | Number of disputes to fetch             |
| `--skip`  | int   | 0       | Records to skip for pagination          |
| `--from`  | int64 | --      | Unix timestamp lower bound (created_at) |
| `--to`    | int64 | --      | Unix timestamp upper bound (created_at) |

Example:

```bash
razorpay disputes list --count 20
```

## Fetch a dispute

```bash
razorpay disputes fetch <dispute_id>
```

Example:

```bash
razorpay disputes fetch disp_AHfqOvkokHwMsN
```

## Accept a dispute

Accepting a dispute means you concede the chargeback. The disputed amount is deducted from your next settlement.

```bash
razorpay disputes accept <dispute_id>
```

Example:

```bash
razorpay disputes accept disp_AHfqOvkokHwMsN
```

> This action is irreversible. Once accepted, the dispute cannot be contested.

## Contest a dispute

Contesting a dispute submits evidence to Razorpay for review. You can save a draft first, then submit when ready.

```bash
razorpay disputes contest <dispute_id> [--action <action>] [--param <key=value>]
```

Flags:

| Flag       | Type   | Description                                        |
| ---------- | ------ | -------------------------------------------------- |
| `--action` | string | `draft` to save without submitting, `submit` to submit for review |
| `--param`  | string | Evidence parameter as `key=value`, repeatable      |

Examples:

```bash
# Save a draft with evidence fields
razorpay disputes contest disp_AHfqOvkokHwMsN \
  --action draft \
  --param "billing_proof=doc_xxx" \
  --param "explanation=Order was delivered on time"

# Submit the contest
razorpay disputes contest disp_AHfqOvkokHwMsN --action submit
```

## Dispute statuses

| Status      | Meaning                                                        |
| ----------- | -------------------------------------------------------------- |
| `open`      | Dispute raised; action required before the deadline            |
| `under_review` | Evidence submitted; awaiting bank decision                  |
| `won`       | Dispute resolved in your favour; no deduction                  |
| `lost`      | Dispute resolved against you; amount deducted                  |
| `closed`    | Dispute closed without a chargeback (e.g. customer retracted)  |
