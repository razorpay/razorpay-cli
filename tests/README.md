# Razorpay CLI — End-to-End Tests

This suite runs the compiled `razorpay` binary as a subprocess against the
real Razorpay **test** API and asserts every command and subcommand behaves
correctly.

The tests are gated by the `e2e` build tag, so a plain `go test ./...` will
never invoke them.

## Prerequisites

Set your Razorpay test API credentials in the environment:

```sh
export RAZORPAY_TEST_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_TEST_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```

For safety, the suite **refuses to run** if `RAZORPAY_TEST_KEY_ID` does not
start with `rzp_test_`. As a fallback, the suite will also accept
`RAZORPAY_KEY_ID` / `RAZORPAY_KEY_SECRET`.

## Running

```sh
go test -tags=e2e -v ./tests/...
```

By default the suite exercises every command except a few that need stateful
prerequisites the test account may not have. Those subtests skip with a clear
hint. To exercise them, set the corresponding env variables before running:

| Command                | Env variable(s) to set                                              |
| ---------------------- | ------------------------------------------------------------------- |
| `payments capture`     | `RAZORPAY_TEST_AUTHORIZED_PAYMENT_ID`, `RAZORPAY_TEST_AUTHORIZED_PAYMENT_AMOUNT` |
| `refunds create`       | `RAZORPAY_TEST_CAPTURED_PAYMENT_ID`                                 |
| `disputes contest`     | `RAZORPAY_TEST_RUN_DESTRUCTIVE=1` (mutates an existing dispute)     |
| `disputes accept`      | `RAZORPAY_TEST_RUN_DESTRUCTIVE=1` (debits the test account)         |

## Coverage

For every command and subcommand the suite verifies:

- `--help` exits cleanly and renders non-empty help text.
- Argument validation (missing flags, bad `--param` format, wrong arg counts,
  unknown commands) produces a non-zero exit with a clear message.
- The success path returns valid JSON whose `entity` field matches the
  resource type returned by the Razorpay API.

The complete list of commands tested lives in `allCommands` in
[`e2e_test.go`](./e2e_test.go). When a new command is added to the CLI, append
it to that list so it picks up `--help` coverage automatically.

## What the tests do to your test account

These tests create real resources on your **test** account:

- One order per run (`orders create`).
- One customer per run (`customers create`, with a unique email).
- A `notes` field update on the most recent payment (if any).

Nothing is deleted or refunded by default. `payments capture`, `refunds create`,
`disputes contest`, and `disputes accept` only run when their dedicated env
variables are set.
