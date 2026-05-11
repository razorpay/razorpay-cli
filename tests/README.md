# Razorpay CLI — End-to-End Tests

This suite runs the compiled `razorpay` binary as a subprocess against the
live Razorpay API and asserts every command behaves correctly. The tests
are gated by the `e2e` build tag so a plain `go test ./...` skips them.

## Prerequisites

Set Razorpay API credentials in the environment. **Any key works** — the
suite does not enforce a `rzp_test_` prefix. Whichever account the key
belongs to is what the tests will create resources on, so most users will
want to point this at a test account:

```sh
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```

If either variable is missing the suite exits with a non-zero status and
prompts you to set them.

## Running

```sh
go test -tags=e2e -v ./tests/...
```

No other env variables are required. Workflows that cannot run on a given
account (e.g. Route account creation when KYC is gated, document uploads
on accounts that disallow them) skip cleanly with the API error inlined.

## Test layout

```
helpers_test.go      build CLI binary, run/runJSON helpers, JSON parsers
meta_test.go         dynamic --help walk, configure flows, validation errors
workflows_test.go    one lifecycle per resource group
```

### Help coverage

`TestHelp` discovers the subcommand tree at runtime by parsing
`--help` output, so newly-added commands pick up `--help` coverage
automatically without any change to the test file.

### Lifecycle workflows

Each `TestXxxLifecycle` exercises one resource group through a homogeneous
flow — **create → fetch → list → update → terminal** — using the ID
returned by the create step to drive every dependent call.

| Workflow                             | Lifecycle covered                                                                |
| ------------------------------------ | -------------------------------------------------------------------------------- |
| `TestOrdersLifecycle`                | create → fetch → list → update → payments                                        |
| `TestCustomersLifecycle`             | create → fetch → list → update                                                   |
| `TestInvoiceItemsLifecycle`          | create → fetch → list → update → delete                                          |
| `TestInvoicesLifecycle`              | create (draft) → fetch → list → update → delete                                  |
| `TestPaymentLinksLifecycle`          | create → fetch → list → update → cancel                                          |
| `TestQRCodesLifecycle`               | create → fetch → list → update → payments → close                                |
| `TestPlansAndSubscriptionsLifecycle` | plan create/fetch/list → subscription create/fetch/list → pause/resume → cancel  |
| `TestDocumentsLifecycle`             | upload (in-memory PNG) → fetch → fetch-content                                   |
| `TestSmartCollectLifecycle`          | customer → virtual account create → fetch → list → payments → add-receiver → close |
| `TestRouteAccountsLifecycle`         | account create → fetch → update → transfer list                                   |

### List-driven workflows

Resources that can only be **created by real-world events** (payments,
refunds, disputes, settlements) discover state from `list` and skip
cleanly when nothing matches:

| Workflow                  | What it covers                                                                       |
| ------------------------- | ------------------------------------------------------------------------------------ |
| `TestPaymentsListDriven`  | list → fetch → card → update notes → capture (when an `authorized` payment exists) → downtime list/fetch |
| `TestRefundsListDriven`   | refund a `captured` payment when one exists; otherwise list → fetch → update         |
| `TestDisputesListDriven`  | list → fetch (accept/contest are deliberately not auto-fired)                        |
| `TestSettlementsListDriven` | list → fetch → recon → instant-list → instant-fetch                               |

## What the tests do to your account

Each run creates a handful of resources: an order, a customer, draft
invoices, payment links, QR codes, a subscription plan + subscription,
and a virtual account. Most are cancelled / closed / deleted in the same
test, but some (orders, customers) are not — Razorpay's API has no delete
operation for them.

Destructive operations on objects the tests did not themselves create —
accepting a real dispute, capturing a real authorised payment — only run
when matching state already exists on the account.
