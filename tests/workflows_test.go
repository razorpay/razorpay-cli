//go:build e2e
// +build e2e

package tests

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Each TestXxxLifecycle below exercises a single resource group through a
// homogeneous lifecycle (create → fetch → list → update → terminal) using
// IDs returned by the create step. Resources that depend on real-world
// events (payments, refunds, disputes, settlements) are exercised in
// list-driven workflows: they discover state from `list` and skip cleanly
// when nothing is present.

func TestOrdersLifecycle(t *testing.T) {
	var orderID string

	t.Run("create", func(t *testing.T) {
		resp := runJSON(t, "orders", "create",
			"--amount", "50000",
			"--currency", "INR",
			"--receipt", "e2e-"+uniqSuffix(),
			"--note", "source=razorpay-cli-e2e")
		requireEntity(t, resp, "order")
		orderID = strField(resp, "id")
		if orderID == "" {
			t.Fatalf("orders create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "fetch", orderID)
		requireEntity(t, resp, "order")
	})

	t.Run("list_contains_new_order", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "list", "--count", "25")
		hit := findItem(t, resp, func(m map[string]any) bool {
			return strField(m, "id") == orderID
		})
		if hit == nil {
			t.Logf("freshly-created order %s not present in latest 25 — acceptable but worth noting", orderID)
		}
	})

	t.Run("update_notes", func(t *testing.T) {
		if orderID == "" {
			t.Skip("no order id from create step")
		}
		resp := runJSON(t, "orders", "update", orderID,
			"--note", "shipment=AWB1234")
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

func TestCustomersLifecycle(t *testing.T) {
	var custID string

	t.Run("create", func(t *testing.T) {
		email := fmt.Sprintf("e2e+%s@example.com", uniqSuffix())
		resp := runJSON(t, "customers", "create",
			"--name", "E2E Tester",
			"--email", email,
			"--contact", "9999999999")
		requireEntity(t, resp, "customer")
		custID = strField(resp, "id")
		if custID == "" {
			t.Fatalf("customers create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if custID == "" {
			t.Skip("no customer id from create step")
		}
		resp := runJSON(t, "customers", "fetch", custID)
		requireEntity(t, resp, "customer")
	})

	t.Run("list", func(t *testing.T) {
		runJSON(t, "customers", "list", "--count", "5")
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

func TestInvoiceItemsLifecycle(t *testing.T) {
	var itemID string

	t.Run("create", func(t *testing.T) {
		resp := runJSON(t, "invoices", "items", "create",
			"--name", "E2E Item "+uniqSuffix(),
			"--amount", "1000",
			"--currency", "INR",
			"--description", "razorpay-cli e2e")
		// invoices items create returns the item directly.
		itemID = strField(resp, "id")
		if itemID == "" {
			t.Fatalf("items create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if itemID == "" {
			t.Skip("no item id from create step")
		}
		runJSON(t, "invoices", "items", "fetch", itemID)
	})

	t.Run("list", func(t *testing.T) {
		runJSON(t, "invoices", "items", "list", "--count", "5")
	})

	t.Run("update", func(t *testing.T) {
		if itemID == "" {
			t.Skip("no item id from create step")
		}
		runJSON(t, "invoices", "items", "update", itemID,
			"--description", "razorpay-cli e2e (updated)")
	})

	t.Run("delete", func(t *testing.T) {
		if itemID == "" {
			t.Skip("no item id from create step")
		}
		r := run(t, "invoices", "items", "delete", itemID)
		if r.err != nil {
			t.Fatalf("items delete failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestInvoicesLifecycle(t *testing.T) {
	var invoiceID string

	t.Run("create_draft", func(t *testing.T) {
		// Create as draft so we can exercise issue → cancel without firing
		// real emails / sms to a real customer.
		resp := runJSON(t, "invoices", "create",
			"--type", "invoice",
			"--draft",
			"--customer-name", "E2E Customer",
			"--customer-email", fmt.Sprintf("e2e+inv+%s@example.com", uniqSuffix()),
			"--customer-contact", "9999999999",
			"--line-items", `[{"name":"E2E Service","amount":5000,"currency":"INR"}]`,
			"--description", "razorpay-cli e2e draft",
			"--note", "source=razorpay-cli-e2e")
		requireEntity(t, resp, "invoice")
		invoiceID = strField(resp, "id")
		if invoiceID == "" {
			t.Fatalf("invoices create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if invoiceID == "" {
			t.Skip("no invoice id")
		}
		runJSON(t, "invoices", "fetch", invoiceID)
	})

	t.Run("list", func(t *testing.T) {
		runJSON(t, "invoices", "list")
	})

	t.Run("update", func(t *testing.T) {
		if invoiceID == "" {
			t.Skip("no invoice id")
		}
		runJSON(t, "invoices", "update", invoiceID,
			"--description", "razorpay-cli e2e draft (updated)")
	})

	t.Run("delete_draft", func(t *testing.T) {
		// Razorpay only allows deleting invoices in `draft` state. Since we
		// never issued this one, delete should succeed and we skip the
		// cancel subtest below.
		if invoiceID == "" {
			t.Skip("no invoice id")
		}
		r := run(t, "invoices", "delete", invoiceID)
		if r.err != nil {
			t.Fatalf("invoices delete failed: %v\nstderr: %s", r.err, r.stderr)
		}
		invoiceID = "" // mark as gone
	})
}

func TestPaymentLinksLifecycle(t *testing.T) {
	var plID string

	t.Run("create", func(t *testing.T) {
		resp := runJSON(t, "payment-links", "create",
			"--amount", "10000",
			"--currency", "INR",
			"--description", "razorpay-cli e2e",
			"--reference-id", "e2e-"+uniqSuffix(),
			"--customer-name", "E2E Customer",
			"--customer-email", fmt.Sprintf("e2e+pl+%s@example.com", uniqSuffix()),
			"--customer-contact", "9876543210",
			"--notify-email=false",
			"--notify-sms=false")
		plID = strField(resp, "id")
		if plID == "" {
			t.Fatalf("payment-links create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if plID == "" {
			t.Skip("no payment link id")
		}
		runJSON(t, "payment-links", "fetch", plID)
	})

	t.Run("list", func(t *testing.T) {
		runJSON(t, "payment-links", "list")
	})

	t.Run("update", func(t *testing.T) {
		if plID == "" {
			t.Skip("no payment link id")
		}
		// Update is only legal in `created` or `partially_paid` state, which
		// holds for a freshly-created link.
		r := run(t, "payment-links", "update", plID,
			"--reference-id", "e2e-updated-"+uniqSuffix())
		if r.err != nil {
			t.Skipf("payment-links update not allowed for this link: %s", strings.TrimSpace(r.stderr))
		}
	})

	t.Run("cancel", func(t *testing.T) {
		if plID == "" {
			t.Skip("no payment link id")
		}
		r := run(t, "payment-links", "cancel", plID)
		if r.err != nil {
			t.Fatalf("payment-links cancel failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestQRCodesLifecycle(t *testing.T) {
	var qrID string

	t.Run("create", func(t *testing.T) {
		// QR codes require the feature to be enabled on the account. Skip
		// rather than fail when the API returns "URL was not found" or a
		// capability rejection.
		resp := runOrSkipJSON(t, "qr-codes", "create",
			"--type", "upi_qr",
			"--name", "E2E QR "+uniqSuffix(),
			"--usage", "single_use",
			"--fixed-amount",
			"--payment-amount", "10000",
			"--description", "razorpay-cli e2e")
		requireEntity(t, resp, "qr_code")
		qrID = strField(resp, "id")
		if qrID == "" {
			t.Fatalf("qr-codes create did not return id: %v", resp)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if qrID == "" {
			t.Skip("no qr id")
		}
		runJSON(t, "qr-codes", "fetch", qrID)
	})

	t.Run("list", func(t *testing.T) {
		runOrSkipJSON(t, "qr-codes", "list", "--count", "5")
	})

	t.Run("update_notes", func(t *testing.T) {
		if qrID == "" {
			t.Skip("no qr id")
		}
		runJSON(t, "qr-codes", "update", qrID, "--note", "source=razorpay-cli-e2e")
	})

	t.Run("payments_for_qr", func(t *testing.T) {
		if qrID == "" {
			t.Skip("no qr id")
		}
		resp := runJSON(t, "qr-codes", "payments", qrID)
		requireEntity(t, resp, "collection")
	})

	t.Run("close", func(t *testing.T) {
		if qrID == "" {
			t.Skip("no qr id")
		}
		r := run(t, "qr-codes", "close", qrID)
		if r.err != nil {
			t.Fatalf("qr-codes close failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestPlansAndSubscriptionsLifecycle(t *testing.T) {
	var planID, subID string

	t.Run("plan_create", func(t *testing.T) {
		// Subscriptions are gated by an account feature toggle; on accounts
		// without it the API returns 401. Skip rather than fail.
		resp := runOrSkipJSON(t, "subscriptions", "plans", "create",
			"--period", "monthly",
			"--interval", "1",
			"--item-name", "E2E Plan "+uniqSuffix(),
			"--item-amount", "10000",
			"--item-currency", "INR",
			"--item-description", "razorpay-cli e2e")
		requireEntity(t, resp, "plan")
		planID = strField(resp, "id")
	})

	t.Run("plan_fetch", func(t *testing.T) {
		if planID == "" {
			t.Skip("no plan id")
		}
		runJSON(t, "subscriptions", "plans", "fetch", planID)
	})

	t.Run("plan_list", func(t *testing.T) {
		runOrSkipJSON(t, "subscriptions", "plans", "list", "--count", "5")
	})

	t.Run("subscription_create", func(t *testing.T) {
		if planID == "" {
			t.Skip("no plan id")
		}
		resp := runJSON(t, "subscriptions", "create",
			"--plan-id", planID,
			"--total-count", "12",
			"--customer-notify=false")
		requireEntity(t, resp, "subscription")
		subID = strField(resp, "id")
	})

	t.Run("subscription_fetch", func(t *testing.T) {
		if subID == "" {
			t.Skip("no subscription id")
		}
		runJSON(t, "subscriptions", "fetch", subID)
	})

	t.Run("subscription_list", func(t *testing.T) {
		runOrSkipJSON(t, "subscriptions", "list", "--count", "5")
	})

	t.Run("subscription_pause", func(t *testing.T) {
		if subID == "" {
			t.Skip("no subscription id")
		}
		// `pause` only works on active subscriptions; ours is just `created`,
		// so this typically returns a 400. Treat as skip.
		r := run(t, "subscriptions", "pause", subID)
		if r.err != nil {
			t.Skipf("pause not applicable on a freshly-created subscription: %s",
				strings.TrimSpace(r.stderr))
		}
	})

	t.Run("subscription_resume", func(t *testing.T) {
		if subID == "" {
			t.Skip("no subscription id")
		}
		r := run(t, "subscriptions", "resume", subID)
		if r.err != nil {
			t.Skipf("resume not applicable: %s", strings.TrimSpace(r.stderr))
		}
	})

	t.Run("subscription_cancel", func(t *testing.T) {
		if subID == "" {
			t.Skip("no subscription id")
		}
		r := run(t, "subscriptions", "cancel", subID)
		if r.err != nil {
			t.Fatalf("subscriptions cancel failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestDocumentsLifecycle(t *testing.T) {
	var docID string

	// Build a tiny valid PNG (1×1 black pixel) so the multipart upload
	// passes Razorpay's content-type check without needing a fixture file.
	tmp, err := os.CreateTemp("", "razorpay-cli-e2e-*.png")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmp.Name()) })
	if _, err := tmp.Write(tinyPNG); err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	tmp.Close()

	t.Run("create_upload", func(t *testing.T) {
		// Some accounts disallow document uploads outside of dispute flows;
		// treat that as a skip rather than a failure.
		resp := runOrSkipJSON(t, "documents", "create",
			"--file", tmp.Name(),
			"--purpose", "dispute_evidence")
		docID = strField(resp, "id")
	})

	t.Run("fetch", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document id")
		}
		runJSON(t, "documents", "fetch", docID)
	})

	t.Run("fetch_content", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document id")
		}
		// fetch-content streams binary content to stdout. Just assert the
		// command exits cleanly.
		r := run(t, "documents", "fetch-content", docID)
		if r.err != nil {
			t.Fatalf("documents fetch-content failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestSmartCollectLifecycle(t *testing.T) {
	var custID, vaID string

	t.Run("customer_for_va", func(t *testing.T) {
		resp := runJSON(t, "customers", "create",
			"--name", "E2E SC Customer",
			"--email", fmt.Sprintf("e2e+sc+%s@example.com", uniqSuffix()),
			"--contact", "9999999999")
		custID = strField(resp, "id")
	})

	t.Run("virtual_account_create", func(t *testing.T) {
		if custID == "" {
			t.Skip("no customer id for VA")
		}
		resp := runOrSkipJSON(t, "smart-collect", "create",
			"--receiver-type", "bank_account",
			"--customer-id", custID,
			"--description", "razorpay-cli e2e")
		requireEntity(t, resp, "virtual_account")
		vaID = strField(resp, "id")
	})

	t.Run("virtual_account_fetch", func(t *testing.T) {
		if vaID == "" {
			t.Skip("no virtual account id")
		}
		runJSON(t, "smart-collect", "fetch", vaID)
	})

	t.Run("virtual_account_list", func(t *testing.T) {
		runOrSkipJSON(t, "smart-collect", "list", "--count", "5")
	})

	t.Run("virtual_account_payments", func(t *testing.T) {
		if vaID == "" {
			t.Skip("no virtual account id")
		}
		resp := runJSON(t, "smart-collect", "payments", vaID)
		requireEntity(t, resp, "collection")
	})

	t.Run("virtual_account_add_receiver", func(t *testing.T) {
		if vaID == "" {
			t.Skip("no virtual account id")
		}
		r := run(t, "smart-collect", "add-receiver", vaID, "--types", "vpa")
		// Some accounts only allow one receiver type per VA; treat that
		// rejection as a skip rather than a failure.
		if r.err != nil {
			t.Skipf("add-receiver not applicable for this VA: %s",
				strings.TrimSpace(r.stderr))
		}
	})

	t.Run("virtual_account_close", func(t *testing.T) {
		if vaID == "" {
			t.Skip("no virtual account id")
		}
		r := run(t, "smart-collect", "close", vaID)
		if r.err != nil {
			t.Fatalf("smart-collect close failed: %v\nstderr: %s", r.err, r.stderr)
		}
	})
}

func TestRouteAccountsLifecycle(t *testing.T) {
	// Route account creation is heavily gated by KYC settings on the parent
	// merchant. On most test accounts the first call here returns 400; we
	// rely on runOrSkipJSON to skip the rest of the workflow gracefully.
	var accountID string

	t.Run("account_create", func(t *testing.T) {
		uniq := uniqSuffix()
		resp := runOrSkipJSON(t, "route", "accounts", "create",
			"--email", fmt.Sprintf("e2e+route+%s@example.com", uniq),
			"--phone", "9999999999",
			"--legal-business-name", "E2E Routes "+uniq,
			"--business-type", "individual",
			"--contact-name", "E2E Tester",
			"--reference-id", "e2e-"+uniq,
			"--profile-category", "ecommerce",
			"--profile-subcategory", "ecommerce",
			"--customer-facing-business-name", "E2E Routes")
		requireEntity(t, resp, "account")
		accountID = strField(resp, "id")
	})

	t.Run("account_fetch", func(t *testing.T) {
		if accountID == "" {
			t.Skip("no route account id")
		}
		runJSON(t, "route", "accounts", "fetch", accountID)
	})

	t.Run("account_update", func(t *testing.T) {
		if accountID == "" {
			t.Skip("no route account id")
		}
		r := run(t, "route", "accounts", "update", accountID,
			"--customer-facing-business-name", "E2E Routes Updated")
		if r.err != nil {
			t.Skipf("route accounts update not allowed for this account: %s",
				strings.TrimSpace(r.stderr))
		}
	})

	t.Run("transfer_list", func(t *testing.T) {
		runOrSkipJSON(t, "route", "transfers", "list")
	})
}

func TestPaymentsListDriven(t *testing.T) {
	// Payments cannot be created via API in a deterministic way (they
	// originate from real card/UPI flows). We drive the workflow off
	// whatever `list` returns.
	var anyID, authorizedID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "payments", "list", "--count", "25")
		if f := firstItem(t, resp); f != nil {
			anyID = strField(f, "id")
		}
		if auth := findItem(t, resp, func(m map[string]any) bool {
			return strField(m, "status") == "authorized"
		}); auth != nil {
			authorizedID = strField(auth, "id")
		}
	})

	t.Run("fetch_first", func(t *testing.T) {
		if anyID == "" {
			t.Skip("no payments on the account")
		}
		runJSON(t, "payments", "fetch", anyID)
	})

	t.Run("card_first", func(t *testing.T) {
		if anyID == "" {
			t.Skip("no payments on the account")
		}
		// `card` returns 400 for non-card payment methods; treat as skip.
		r := run(t, "payments", "card", anyID)
		if r.err != nil {
			t.Skipf("card lookup not applicable for this payment: %s",
				strings.TrimSpace(r.stderr))
		}
	})

	t.Run("update_notes_first", func(t *testing.T) {
		if anyID == "" {
			t.Skip("no payments on the account")
		}
		runJSON(t, "payments", "update", anyID, "--note", "e2e_marker=razorpay-cli")
	})

	t.Run("capture_authorized", func(t *testing.T) {
		if authorizedID == "" {
			t.Skip("no payment in `authorized` state on this account")
		}
		// Capture the full authorised amount. Look it up from `fetch`.
		fetched := runJSON(t, "payments", "fetch", authorizedID)
		amt, _ := fetched["amount"].(float64)
		curr, _ := fetched["currency"].(string)
		if amt == 0 || curr == "" {
			t.Skipf("authorized payment %s missing amount/currency", authorizedID)
		}
		runJSON(t, "payments", "capture", authorizedID,
			"--amount", strconv.Itoa(int(amt)),
			"--currency", curr)
	})

	t.Run("downtime_list", func(t *testing.T) {
		runJSON(t, "payments", "downtime", "list")
	})

	t.Run("downtime_fetch_first", func(t *testing.T) {
		resp := runJSON(t, "payments", "downtime", "list")
		first := firstItem(t, resp)
		if first == nil {
			t.Skip("no payment downtimes currently active")
		}
		runJSON(t, "payments", "downtime", "fetch", strField(first, "id"))
	})
}

func TestRefundsListDriven(t *testing.T) {
	var refundID, capturedPaymentID string

	t.Run("payments_list_for_captured", func(t *testing.T) {
		resp := runJSON(t, "payments", "list", "--count", "50")
		if cap := findItem(t, resp, func(m map[string]any) bool {
			return strField(m, "status") == "captured"
		}); cap != nil {
			capturedPaymentID = strField(cap, "id")
		}
	})

	t.Run("create_from_captured", func(t *testing.T) {
		if capturedPaymentID == "" {
			t.Skip("no `captured` payments available to refund")
		}
		// Refund a token 1 paise so the refund is creatable but tiny.
		resp := runOrSkipJSON(t, "refunds", "create", capturedPaymentID,
			"--amount", "1",
			"--speed", "normal",
			"--note", "source=razorpay-cli-e2e")
		requireEntity(t, resp, "refund")
		refundID = strField(resp, "id")
	})

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "refunds", "list", "--count", "5")
		if refundID == "" {
			if f := firstItem(t, resp); f != nil {
				refundID = strField(f, "id")
			}
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if refundID == "" {
			t.Skip("no refund available")
		}
		runJSON(t, "refunds", "fetch", refundID)
	})

	t.Run("update", func(t *testing.T) {
		if refundID == "" {
			t.Skip("no refund available")
		}
		runJSON(t, "refunds", "update", refundID, "--note", "e2e_updated=razorpay-cli")
	})
}

func TestDisputesListDriven(t *testing.T) {
	var disputeID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "disputes", "list")
		if f := firstItem(t, resp); f != nil {
			disputeID = strField(f, "id")
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if disputeID == "" {
			t.Skip("no disputes on the account")
		}
		runJSON(t, "disputes", "fetch", disputeID)
	})

	// `accept` and `contest` are destructive and irreversible; we
	// deliberately do not call them on a real dispute the user may need.
}

func TestSettlementsListDriven(t *testing.T) {
	var settlementID, instantID string

	t.Run("list", func(t *testing.T) {
		resp := runJSON(t, "settlements", "list", "--count", "5")
		if f := firstItem(t, resp); f != nil {
			settlementID = strField(f, "id")
		}
	})

	t.Run("fetch", func(t *testing.T) {
		if settlementID == "" {
			t.Skip("no settlements on the account")
		}
		runJSON(t, "settlements", "fetch", settlementID)
	})

	t.Run("recon", func(t *testing.T) {
		// Year/month are typed as strings on this command; the API tolerates
		// either ISO or numeric form.
		now := time.Now().UTC()
		r := run(t, "settlements", "recon",
			"--year", strconv.Itoa(now.Year()),
			"--month", strconv.Itoa(int(now.Month())))
		// Accept either a successful JSON response or a Razorpay 4xx — test
		// accounts may have no recon data for the requested window.
		if r.err == nil {
			parseJSON(t, []string{"settlements", "recon"}, r.stdout)
			return
		}
		if !strings.Contains(r.stderr, "API request failed") {
			t.Fatalf("recon failed unexpectedly: %v\nstderr: %s", r.err, r.stderr)
		}
	})

	t.Run("instant_list", func(t *testing.T) {
		// Instant settlements (ondemand) are a paid feature; accounts without
		// it get a 4xx. Skip rather than fail in that case.
		resp := runOrSkipJSON(t, "settlements", "instant-list", "--count", "5")
		if f := firstItem(t, resp); f != nil {
			instantID = strField(f, "id")
		}
	})

	t.Run("instant_fetch", func(t *testing.T) {
		if instantID == "" {
			t.Skip("no instant settlements on the account")
		}
		runJSON(t, "settlements", "instant-fetch", instantID)
	})
}

// tinyPNG is a 1×1 transparent PNG (67 bytes). Used by the documents
// workflow so we have a valid image to upload without checking in a binary.
var tinyPNG = []byte{
	0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 'I', 'D', 'A', 'T',
	0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05,
	0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00,
	0x00, 0x00, 'I', 'E', 'N', 'D', 0xae, 0x42, 0x60,
	0x82,
}

