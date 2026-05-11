//go:build e2e
// +build e2e

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// allCommands lists every command and subcommand registered on the CLI.
// Adding a new command here ensures it gets a --help smoke test.
var allCommands = []string{
	"",
	"configure",
	"payments", "payments list", "payments fetch", "payments capture", "payments update", "payments transfers",
	"orders", "orders list", "orders fetch", "orders create", "orders update", "orders payments",
	"customers", "customers list", "customers fetch", "customers create", "customers update",
	"refunds", "refunds list", "refunds fetch", "refunds create", "refunds update",
	"settlements", "settlements list", "settlements fetch", "settlements recon",
	"disputes", "disputes list", "disputes fetch", "disputes accept", "disputes contest",
}

// TestHelp asserts every command and subcommand renders --help without error.
func TestHelp(t *testing.T) {
	for _, c := range allCommands {
		c := c
		name := "root"
		if c != "" {
			name = strings.ReplaceAll(c, " ", "_")
		}
		t.Run(name, func(t *testing.T) {
			args := append(strings.Fields(c), "--help")
			r := run(t, args...)
			if r.err != nil {
				t.Fatalf("--help failed: %v\nstderr: %s", r.err, r.stderr)
			}
			if strings.TrimSpace(r.stdout) == "" {
				t.Fatalf("--help produced no stdout")
			}
		})
	}
}

// TestValidationErrors exercises the CLI's argument-validation paths so we
// catch regressions in required flags, arg counts, and parameter parsing.
func TestValidationErrors(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		expect string
	}{
		{"payments_fetch_missing_id", []string{"payments", "fetch"}, "accepts 1 arg"},
		{"payments_capture_missing_amount", []string{"payments", "capture", "pay_dummy"}, `required flag(s) "amount"`},
		{"orders_create_missing_amount", []string{"orders", "create"}, `required flag(s) "amount"`},
		{"refunds_create_missing_id", []string{"refunds", "create"}, "accepts 1 arg"},
		{"orders_update_bad_param", []string{"orders", "update", "order_x", "--param", "noequals"}, "expected format key=value"},
		{"unknown_command", []string{"bogus"}, "unknown command"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := run(t, c.args...)
			if r.err == nil {
				t.Fatalf("expected non-zero exit, got success\nstdout: %s", r.stdout)
			}
			combined := r.stdout + r.stderr
			if !strings.Contains(combined, c.expect) {
				t.Fatalf("expected output to contain %q, got:\n%s", c.expect, combined)
			}
		})
	}
}

// TestConfigure covers the four ways credentials can be supplied (both flags,
// both stdin, mixed) and the empty-input error case.
func TestConfigure(t *testing.T) {
	assertSaved := func(t *testing.T, r result) {
		t.Helper()
		if r.err != nil {
			t.Fatalf("configure failed: %v\nstdout: %s\nstderr: %s", r.err, r.stdout, r.stderr)
		}
		if !strings.Contains(r.stdout, "Credentials saved to") {
			t.Fatalf("expected success message, got:\n%s", r.stdout)
		}
	}

	t.Run("both_via_flags", func(t *testing.T) {
		r := run(t, "configure", "--key-id", "rzp_test_flagID", "--key-secret", "flag_secret")
		assertSaved(t, r)
	})

	t.Run("both_via_stdin", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("rzp_test_stdinID\nstdin_secret\n"), "configure")
		assertSaved(t, r)
	})

	t.Run("flag_id_stdin_secret", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("only_secret\n"),
			"configure", "--key-id", "rzp_test_idflag")
		assertSaved(t, r)
	})

	t.Run("flag_secret_stdin_id", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("rzp_test_idstdin\n"),
			"configure", "--key-secret", "only_secret_flag")
		assertSaved(t, r)
	})

	t.Run("empty_input_errors", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("\n\n"), "configure")
		if r.err == nil {
			t.Fatalf("expected error for empty input, got success\nstdout: %s", r.stdout)
		}
		if !strings.Contains(r.stdout+r.stderr, "cannot be empty") {
			t.Fatalf("expected 'cannot be empty' error, got: %s", r.stdout+r.stderr)
		}
	})
}

// TestPayments exercises every `payments` subcommand. fetch/update/transfers
// reuse the first ID returned by list. capture requires an authorized payment
// — set RAZORPAY_TEST_AUTHORIZED_PAYMENT_ID and ..._AMOUNT to exercise it.
func TestPayments(t *testing.T) {
	var firstID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "payments", "list", "--count", "5")
		requireEntity(t, resp, "collection")
		if f := firstItem(t, resp); f != nil {
			firstID, _ = f["id"].(string)
		}
	})

	t.Run("list_with_window", func(t *testing.T) {
		// Use a wide window so the call always succeeds, regardless of activity.
		resp := runJSON(t, "payments", "list",
			"--count", "5",
			"--skip", "0",
			"--from", "1",
			"--to", strconv.FormatInt(time.Now().Unix(), 10))
		requireEntity(t, resp, "collection")
	})

	t.Run("fetch", func(t *testing.T) {
		if firstID == "" {
			t.Skip("no payments available on the test account")
		}
		resp := runJSON(t, "payments", "fetch", firstID)
		requireEntity(t, resp, "payment")
	})

	t.Run("transfers", func(t *testing.T) {
		if firstID == "" {
			t.Skip("no payments available on the test account")
		}
		resp := runJSON(t, "payments", "transfers", firstID)
		requireEntity(t, resp, "collection")
	})

	t.Run("update_notes", func(t *testing.T) {
		if firstID == "" {
			t.Skip("no payments available on the test account")
		}
		resp := runJSON(t, "payments", "update", firstID,
			"--param", "notes[e2e_marker]=razorpay-cli")
		requireEntity(t, resp, "payment")
	})

	t.Run("capture", func(t *testing.T) {
		payID := os.Getenv("RAZORPAY_TEST_AUTHORIZED_PAYMENT_ID")
		amount := os.Getenv("RAZORPAY_TEST_AUTHORIZED_PAYMENT_AMOUNT")
		if payID == "" || amount == "" {
			t.Skip("set RAZORPAY_TEST_AUTHORIZED_PAYMENT_ID and RAZORPAY_TEST_AUTHORIZED_PAYMENT_AMOUNT to test capture")
		}
		resp := runJSON(t, "payments", "capture", payID,
			"--amount", amount,
			"--currency", "INR")
		requireEntity(t, resp, "payment")
	})
}

// TestOrders covers the full orders surface. We create a fresh order so
// fetch/update/payments have a deterministic target.
func TestOrders(t *testing.T) {
	var orderID string

	t.Run("create", func(t *testing.T) {
		receipt := fmt.Sprintf("e2e-%d", time.Now().UnixNano())
		resp := runJSON(t, "orders", "create",
			"--amount", "50000",
			"--currency", "INR",
			"--receipt", receipt,
			"--param", "notes[source]=razorpay-cli-e2e")
		requireEntity(t, resp, "order")
		orderID, _ = resp["id"].(string)
		if orderID == "" {
			t.Fatalf("orders create did not return id: %v", resp)
		}
	})

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "orders", "list", "--count", "5")
		requireEntity(t, resp, "collection")
	})

	t.Run("list_filtered_by_status", func(t *testing.T) {
		resp := runJSON(t, "orders", "list", "--status", "created", "--count", "5")
		requireEntity(t, resp, "collection")
	})

	t.Run("fetch", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "fetch", orderID)
		requireEntity(t, resp, "order")
	})

	t.Run("update", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "update", orderID,
			"--param", "notes[shipment]=AWB1234")
		requireEntity(t, resp, "order")
	})

	t.Run("payments_for_order", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "payments", orderID)
		requireEntity(t, resp, "collection")
	})
}

// TestCustomers covers every `customers` subcommand. A fresh customer is
// created so fetch/update target a known record.
func TestCustomers(t *testing.T) {
	var custID string

	t.Run("create", func(t *testing.T) {
		// Email must be unique per merchant, so include a nano timestamp.
		email := fmt.Sprintf("e2e+%d@example.com", time.Now().UnixNano())
		resp := runJSON(t, "customers", "create",
			"--name", "E2E Tester",
			"--email", email,
			"--contact", "9999999999")
		requireEntity(t, resp, "customer")
		custID, _ = resp["id"].(string)
		if custID == "" {
			t.Fatalf("customers create did not return id: %v", resp)
		}
	})

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "customers", "list", "--count", "5")
		requireEntity(t, resp, "collection")
	})

	t.Run("fetch", func(t *testing.T) {
		if custID == "" {
			t.Skip("no customer id from create step")
		}
		resp := runJSON(t, "customers", "fetch", custID)
		requireEntity(t, resp, "customer")
	})

	t.Run("update", func(t *testing.T) {
		if custID == "" {
			t.Skip("no customer id from create step")
		}
		resp := runJSON(t, "customers", "update", custID,
			"--name", "E2E Tester Updated")
		requireEntity(t, resp, "customer")
	})
}

// TestRefunds covers every `refunds` subcommand. `create` needs a captured
// payment id (which the test account must already own) — set
// RAZORPAY_TEST_CAPTURED_PAYMENT_ID to exercise it. fetch/update fall back
// to the first refund from `list` when present.
func TestRefunds(t *testing.T) {
	var refundID string

	t.Run("create", func(t *testing.T) {
		payID := os.Getenv("RAZORPAY_TEST_CAPTURED_PAYMENT_ID")
		if payID == "" {
			t.Skip("set RAZORPAY_TEST_CAPTURED_PAYMENT_ID to exercise refunds create")
		}
		resp := runJSON(t, "refunds", "create", payID,
			"--speed", "normal",
			"--param", "notes[e2e]=razorpay-cli")
		requireEntity(t, resp, "refund")
		refundID, _ = resp["id"].(string)
	})

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "refunds", "list", "--count", "5")
		requireEntity(t, resp, "collection")
		if refundID == "" {
			if f := firstItem(t, resp); f != nil {
				refundID, _ = f["id"].(string)
			}
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if refundID == "" {
			t.Skip("no refund available on the test account")
		}
		resp := runJSON(t, "refunds", "fetch", refundID)
		requireEntity(t, resp, "refund")
	})

	t.Run("update", func(t *testing.T) {
		if refundID == "" {
			t.Skip("no refund available on the test account")
		}
		resp := runJSON(t, "refunds", "update", refundID,
			"--param", "notes[e2e_updated]=razorpay-cli")
		requireEntity(t, resp, "refund")
	})
}

// TestSettlements covers every `settlements` subcommand. `fetch` is skipped
// when the test account has no settlements; `recon` accepts either a
// successful collection response or a Razorpay 4xx (test accounts may have
// no recon data for the requested window).
func TestSettlements(t *testing.T) {
	var settlementID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "settlements", "list", "--count", "5")
		requireEntity(t, resp, "collection")
		if f := firstItem(t, resp); f != nil {
			settlementID, _ = f["id"].(string)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if settlementID == "" {
			t.Skip("no settlements on the test account")
		}
		resp := runJSON(t, "settlements", "fetch", settlementID)
		requireEntity(t, resp, "settlement")
	})

	t.Run("recon", func(t *testing.T) {
		now := time.Now().UTC()
		r := run(t, "settlements", "recon",
			"--year", strconv.Itoa(now.Year()),
			"--month", strconv.Itoa(int(now.Month())))
		if r.err == nil {
			// Success path: must be valid JSON.
			var v map[string]any
			if err := json.Unmarshal([]byte(strings.TrimSpace(r.stdout)), &v); err != nil {
				t.Fatalf("recon returned non-JSON output: %v\nstdout: %s", err, r.stdout)
			}
			return
		}
		// Non-zero exit is acceptable only if it was an API 4xx (i.e. the
		// CLI did its job and surfaced a Razorpay error).
		if !strings.Contains(r.stderr, "API request failed") {
			t.Fatalf("recon failed with unexpected error: %v\nstderr: %s\nstdout: %s",
				r.err, r.stderr, r.stdout)
		}
	})
}

// TestDisputes covers every `disputes` subcommand. fetch/contest/accept all
// require a real dispute; accept and contest are destructive and gated behind
// RAZORPAY_TEST_RUN_DESTRUCTIVE=1.
func TestDisputes(t *testing.T) {
	var disputeID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "disputes", "list", "--count", "5")
		requireEntity(t, resp, "collection")
		if f := firstItem(t, resp); f != nil {
			disputeID, _ = f["id"].(string)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if disputeID == "" {
			t.Skip("no disputes on the test account")
		}
		resp := runJSON(t, "disputes", "fetch", disputeID)
		requireEntity(t, resp, "dispute")
	})

	t.Run("contest_draft", func(t *testing.T) {
		if disputeID == "" {
			t.Skip("no disputes on the test account")
		}
		if os.Getenv("RAZORPAY_TEST_RUN_DESTRUCTIVE") != "1" {
			t.Skip("set RAZORPAY_TEST_RUN_DESTRUCTIVE=1 to mutate dispute state")
		}
		resp := runJSON(t, "disputes", "contest", disputeID,
			"--action", "draft",
			"--param", "evidence[summary]=razorpay-cli e2e draft")
		requireEntity(t, resp, "dispute")
	})

	t.Run("accept", func(t *testing.T) {
		if disputeID == "" {
			t.Skip("no disputes on the test account")
		}
		if os.Getenv("RAZORPAY_TEST_RUN_DESTRUCTIVE") != "1" {
			t.Skip("set RAZORPAY_TEST_RUN_DESTRUCTIVE=1 to accept a dispute (debits the test account)")
		}
		resp := runJSON(t, "disputes", "accept", disputeID)
		requireEntity(t, resp, "dispute")
	})
}
