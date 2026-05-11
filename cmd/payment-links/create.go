package paymentlinks

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new payment link",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		// ── Core fields ──────────────────────────────────────────────────────
		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		description, _ := cmd.Flags().GetString("description")
		referenceID, _ := cmd.Flags().GetString("reference-id")
		acceptPartial, _ := cmd.Flags().GetBool("accept-partial")
		firstMinPartialAmount, _ := cmd.Flags().GetInt64("first-min-partial-amount")
		upiLink, _ := cmd.Flags().GetBool("upi-link")
		expireBy, _ := cmd.Flags().GetInt64("expire-by")
		callbackURL, _ := cmd.Flags().GetString("callback-url")
		callbackMethod, _ := cmd.Flags().GetString("callback-method")
		reminderEnable, _ := cmd.Flags().GetBool("reminder-enable")
		customerName, _ := cmd.Flags().GetString("customer-name")
		customerContact, _ := cmd.Flags().GetString("customer-contact")
		customerEmail, _ := cmd.Flags().GetString("customer-email")
		notifySMS, _ := cmd.Flags().GetBool("notify-sms")
		notifyEmail, _ := cmd.Flags().GetBool("notify-email")
		notes, _ := cmd.Flags().GetStringArray("note")

		// ── options.checkout ─────────────────────────────────────────────────
		checkoutName, _ := cmd.Flags().GetString("checkout-name")
		checkoutHideTopbar, _ := cmd.Flags().GetBool("checkout-hide-topbar")
		checkoutConfig, _ := cmd.Flags().GetString("checkout-config")
		methodNetbanking, _ := cmd.Flags().GetBool("method-netbanking")
		methodCard, _ := cmd.Flags().GetBool("method-card")
		methodUPI, _ := cmd.Flags().GetBool("method-upi")
		methodWallet, _ := cmd.Flags().GetBool("method-wallet")
		labelMinAmount, _ := cmd.Flags().GetString("label-min-amount")
		labelPartialAmount, _ := cmd.Flags().GetString("label-partial-amount")
		labelPartialDesc, _ := cmd.Flags().GetString("label-partial-description")
		labelFullAmount, _ := cmd.Flags().GetString("label-full-amount")

		// ── options.hosted_page ───────────────────────────────────────────────
		hostedLabelReceipt, _ := cmd.Flags().GetString("hosted-label-receipt")
		hostedLabelDesc, _ := cmd.Flags().GetString("hosted-label-description")
		hostedLabelAmountPayable, _ := cmd.Flags().GetString("hosted-label-amount-payable")
		hostedLabelAmountPaid, _ := cmd.Flags().GetString("hosted-label-amount-paid")
		hostedLabelPartialDue, _ := cmd.Flags().GetString("hosted-label-partial-amount-due")
		hostedLabelPartialPaid, _ := cmd.Flags().GetString("hosted-label-partial-amount-paid")
		hostedLabelExpireBy, _ := cmd.Flags().GetString("hosted-label-expire-by")
		hostedLabelExpiredOn, _ := cmd.Flags().GetString("hosted-label-expired-on")
		hostedLabelAmountDue, _ := cmd.Flags().GetString("hosted-label-amount-due")
		hostedShowIssuedTo, _ := cmd.Flags().GetBool("hosted-show-issued-to")

		// ── options.order ────────────────────────────────────────────────────
		offers, _ := cmd.Flags().GetStringArray("offer")
		transfersJSON, _ := cmd.Flags().GetString("transfers")
		orderMethod, _ := cmd.Flags().GetString("order-method")
		bankAccountNumber, _ := cmd.Flags().GetString("bank-account-number")
		bankAccountName, _ := cmd.Flags().GetString("bank-account-name")
		bankAccountIFSC, _ := cmd.Flags().GetString("bank-account-ifsc")

		// ── Validation ────────────────────────────────────────────────────────
		if amount <= 0 {
			return fmt.Errorf("--amount is required and must be > 0")
		}
		if acceptPartial && firstMinPartialAmount <= 0 {
			return fmt.Errorf("--first-min-partial-amount is required when --accept-partial is set")
		}
		if orderMethod != "" && bankAccountNumber == "" {
			return fmt.Errorf("--bank-account-number, --bank-account-name and --bank-account-ifsc are required with --order-method")
		}

		// ── Core body ─────────────────────────────────────────────────────────
		body := map[string]any{
			"amount":   amount,
			"currency": currency,
		}
		if description != "" {
			body["description"] = description
		}
		if referenceID != "" {
			body["reference_id"] = referenceID
		}
		if acceptPartial {
			body["accept_partial"] = true
			body["first_min_partial_amount"] = firstMinPartialAmount
		}
		if upiLink {
			body["upi_link"] = true
		}
		if expireBy > 0 {
			body["expire_by"] = expireBy
		}
		if callbackURL != "" {
			body["callback_url"] = callbackURL
			if callbackMethod != "" {
				body["callback_method"] = callbackMethod
			}
		}
		if cmd.Flags().Changed("reminder-enable") {
			body["reminder_enable"] = reminderEnable
		}

		// ── customer ──────────────────────────────────────────────────────────
		customer := map[string]any{}
		if customerName != "" {
			customer["name"] = customerName
		}
		if customerContact != "" {
			customer["contact"] = customerContact
		}
		if customerEmail != "" {
			customer["email"] = customerEmail
		}
		if len(customer) > 0 {
			body["customer"] = customer
		}

		// ── notify ────────────────────────────────────────────────────────────
		if cmd.Flags().Changed("notify-sms") || cmd.Flags().Changed("notify-email") {
			body["notify"] = map[string]any{
				"sms":   notifySMS,
				"email": notifyEmail,
			}
		}

		// ── notes ─────────────────────────────────────────────────────────────
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		// ── options ───────────────────────────────────────────────────────────
		options := map[string]any{}

		// options.checkout
		checkout := map[string]any{}
		if checkoutName != "" {
			checkout["name"] = checkoutName
		}
		if cmd.Flags().Changed("checkout-hide-topbar") {
			checkout["theme"] = map[string]any{"hide_topbar": checkoutHideTopbar}
		}
		// payment methods
		method := map[string]any{}
		if cmd.Flags().Changed("method-netbanking") {
			method["netbanking"] = methodNetbanking
		}
		if cmd.Flags().Changed("method-card") {
			method["card"] = methodCard
		}
		if cmd.Flags().Changed("method-upi") {
			method["upi"] = methodUPI
		}
		if cmd.Flags().Changed("method-wallet") {
			method["wallet"] = methodWallet
		}
		if len(method) > 0 {
			checkout["method"] = method
		}
		// checkout config (raw JSON)
		if checkoutConfig != "" {
			var configObj any
			if err := json.Unmarshal([]byte(checkoutConfig), &configObj); err != nil {
				return fmt.Errorf("--checkout-config is not valid JSON: %w", err)
			}
			checkout["config"] = configObj
		}
		// partial payment labels
		partialPayment := map[string]any{}
		if labelMinAmount != "" {
			partialPayment["min_amount_label"] = labelMinAmount
		}
		if labelPartialAmount != "" {
			partialPayment["partial_amount_label"] = labelPartialAmount
		}
		if labelPartialDesc != "" {
			partialPayment["partial_amount_description"] = labelPartialDesc
		}
		if labelFullAmount != "" {
			partialPayment["full_amount_label"] = labelFullAmount
		}
		if len(partialPayment) > 0 {
			checkout["partial_payment"] = partialPayment
		}
		if len(checkout) > 0 {
			options["checkout"] = checkout
		}

		// options.hosted_page
		hostedLabel := map[string]any{}
		if hostedLabelReceipt != "" {
			hostedLabel["receipt"] = hostedLabelReceipt
		}
		if hostedLabelDesc != "" {
			hostedLabel["description"] = hostedLabelDesc
		}
		if hostedLabelAmountPayable != "" {
			hostedLabel["amount_payable"] = hostedLabelAmountPayable
		}
		if hostedLabelAmountPaid != "" {
			hostedLabel["amount_paid"] = hostedLabelAmountPaid
		}
		if hostedLabelPartialDue != "" {
			hostedLabel["partial_amount_due"] = hostedLabelPartialDue
		}
		if hostedLabelPartialPaid != "" {
			hostedLabel["partial_amount_paid"] = hostedLabelPartialPaid
		}
		if hostedLabelExpireBy != "" {
			hostedLabel["expire_by"] = hostedLabelExpireBy
		}
		if hostedLabelExpiredOn != "" {
			hostedLabel["expired_on"] = hostedLabelExpiredOn
		}
		if hostedLabelAmountDue != "" {
			hostedLabel["amount_due"] = hostedLabelAmountDue
		}
		hostedPage := map[string]any{}
		if len(hostedLabel) > 0 {
			hostedPage["label"] = hostedLabel
		}
		if cmd.Flags().Changed("hosted-show-issued-to") {
			hostedPage["show_preferences"] = map[string]any{"issued_to": hostedShowIssuedTo}
		}
		if len(hostedPage) > 0 {
			options["hosted_page"] = hostedPage
		}

		// options.order
		order := map[string]any{}
		if len(offers) > 0 {
			order["offers"] = offers
		}
		if transfersJSON != "" {
			var transfersObj any
			if err := json.Unmarshal([]byte(transfersJSON), &transfersObj); err != nil {
				return fmt.Errorf("--transfers is not valid JSON: %w", err)
			}
			order["transfers"] = transfersObj
		}
		if orderMethod != "" {
			order["method"] = orderMethod
			order["bank_account"] = map[string]any{
				"account_number": bankAccountNumber,
				"name":           bankAccountName,
				"ifsc":           bankAccountIFSC,
			}
		}
		if len(order) > 0 {
			options["order"] = order
		}

		if len(options) > 0 {
			body["options"] = options
		}

		data, err := client.Post(basePath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(createCmd)

	// Core
	createCmd.Flags().Int64("amount", 0, "Amount in paise (required, e.g. 30000 for ₹300)")
	createCmd.Flags().String("currency", "INR", "Currency code (default INR)")
	createCmd.Flags().String("description", "", "Payment description (max 2048 characters)")
	createCmd.Flags().String("reference-id", "", "Unique reference ID for internal tracking (max 40 characters)")
	createCmd.Flags().Bool("accept-partial", false, "Allow customer to pay in partial amounts")
	createCmd.Flags().Int64("first-min-partial-amount", 0, "Minimum first partial payment amount in paise (required with --accept-partial)")
	createCmd.Flags().Bool("upi-link", false, "Create a UPI payment link instead of a standard link")
	createCmd.Flags().Int64("expire-by", 0, "Link expiry as Unix timestamp (max 6 months from now)")
	createCmd.Flags().String("callback-url", "", "Redirect URL after payment completion")
	createCmd.Flags().String("callback-method", "get", "Callback HTTP method (must be get)")
	createCmd.Flags().Bool("reminder-enable", true, "Send payment reminders to customer")

	// Customer
	createCmd.Flags().String("customer-name", "", "Customer name")
	createCmd.Flags().String("customer-contact", "", "Customer phone number")
	createCmd.Flags().String("customer-email", "", "Customer email address")

	// Notify
	createCmd.Flags().Bool("notify-sms", true, "Send SMS notification to customer")
	createCmd.Flags().Bool("notify-email", true, "Send email notification to customer")

	// Notes
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")

	// options.checkout — business name & theme
	createCmd.Flags().String("checkout-name", "", "Business name to display on the checkout page")
	createCmd.Flags().Bool("checkout-hide-topbar", false, "Hide the top bar on the checkout page")

	// options.checkout — payment methods
	createCmd.Flags().Bool("method-netbanking", false, "Enable/disable netbanking as a payment method")
	createCmd.Flags().Bool("method-card", false, "Enable/disable card as a payment method")
	createCmd.Flags().Bool("method-upi", false, "Enable/disable UPI as a payment method")
	createCmd.Flags().Bool("method-wallet", false, "Enable/disable wallet as a payment method")
	for _, name := range []string{"method-netbanking", "method-card", "method-upi", "method-wallet"} {
		createCmd.Flags().Lookup(name).NoOptDefVal = ""
	}

	// options.checkout — config (complex JSON)
	createCmd.Flags().String("checkout-config", "", "Checkout display config as raw JSON (for advanced block/sequence customisation)")

	// options.checkout — partial payment labels
	createCmd.Flags().String("label-min-amount", "", "Label for the minimum partial payment amount field")
	createCmd.Flags().String("label-partial-amount", "", "Label for the partial payment amount field")
	createCmd.Flags().String("label-partial-description", "", "Description text below the partial payment field")
	createCmd.Flags().String("label-full-amount", "", "Label for the full payment amount option")

	// options.hosted_page — payment details labels
	createCmd.Flags().String("hosted-label-receipt", "", "Label for receipt field on hosted page")
	createCmd.Flags().String("hosted-label-description", "", "Label for description field on hosted page")
	createCmd.Flags().String("hosted-label-amount-payable", "", "Label for amount payable field on hosted page")
	createCmd.Flags().String("hosted-label-amount-paid", "", "Label for amount paid field on hosted page")
	createCmd.Flags().String("hosted-label-partial-amount-due", "", "Label for partial amount due field on hosted page")
	createCmd.Flags().String("hosted-label-partial-amount-paid", "", "Label for partial amount paid field on hosted page")
	createCmd.Flags().String("hosted-label-expire-by", "", "Label for expiry date field on hosted page")
	createCmd.Flags().String("hosted-label-expired-on", "", "Text shown when the link has expired on hosted page")
	createCmd.Flags().String("hosted-label-amount-due", "", "Label for amount due field on hosted page")
	createCmd.Flags().Bool("hosted-show-issued-to", true, "Show/hide the issued-to field on hosted page")

	// options.order — offers
	createCmd.Flags().StringArray("offer", nil, "Offer ID to apply (repeatable, e.g. --offer offer_ABC --offer offer_XYZ)")

	// options.order — transfers (complex JSON)
	createCmd.Flags().String("transfers", "", "Route transfers as raw JSON array (e.g. '[{\"account\":\"acc_XX\",\"amount\":500,\"currency\":\"INR\"}]')")

	// options.order — third-party validation
	createCmd.Flags().String("order-method", "", "Payment method for third-party validation (e.g. netbanking)")
	createCmd.Flags().String("bank-account-number", "", "Bank account number for third-party validation")
	createCmd.Flags().String("bank-account-name", "", "Account holder name for third-party validation")
	createCmd.Flags().String("bank-account-ifsc", "", "Bank IFSC code for third-party validation")
}
