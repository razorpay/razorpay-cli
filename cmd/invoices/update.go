package invoices

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <invoice_id>",
	Short: "Update an invoice",
	Long: `Update an invoice.

Updatable in ALL states (draft and issued):
  --description, --expire-by, --note

Updatable in DRAFT state only:
  --customer-id, --customer-*, --billing-*, --shipping-*, --line-items,
  --currency, --date, --terms, --sms-notify, --email-notify, --partial-payment`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		description, _ := cmd.Flags().GetString("description")
		expireBy, _ := cmd.Flags().GetInt64("expire-by")
		notes, _ := cmd.Flags().GetStringArray("note")
		// draft-only
		customerID, _ := cmd.Flags().GetString("customer-id")
		customerName, _ := cmd.Flags().GetString("customer-name")
		customerContact, _ := cmd.Flags().GetString("customer-contact")
		customerEmail, _ := cmd.Flags().GetString("customer-email")
		lineItemsJSON, _ := cmd.Flags().GetString("line-items")
		smsNotify, _ := cmd.Flags().GetBool("sms-notify")
		emailNotify, _ := cmd.Flags().GetBool("email-notify")
		partialPayment, _ := cmd.Flags().GetBool("partial-payment")
		currency, _ := cmd.Flags().GetString("currency")
		date, _ := cmd.Flags().GetInt64("date")
		terms, _ := cmd.Flags().GetString("terms")

		body := map[string]any{}

		// All states
		if description != "" {
			body["description"] = description
		}
		if expireBy > 0 {
			body["expire_by"] = expireBy
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		// Draft-only
		if customerID != "" {
			body["customer_id"] = customerID
		}
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
		if billing := buildAddress(cmd, "billing"); len(billing) > 0 {
			customer["billing_address"] = billing
		}
		if shipping := buildAddress(cmd, "shipping"); len(shipping) > 0 {
			customer["shipping_address"] = shipping
		}
		if len(customer) > 0 {
			body["customer"] = customer
		}
		if lineItemsJSON != "" {
			var items any
			if err := json.Unmarshal([]byte(lineItemsJSON), &items); err != nil {
				return fmt.Errorf("--line-items is not valid JSON: %w", err)
			}
			body["line_items"] = items
		}
		if cmd.Flags().Changed("sms-notify") {
			body["sms_notify"] = smsNotify
		}
		if cmd.Flags().Changed("email-notify") {
			body["email_notify"] = emailNotify
		}
		if cmd.Flags().Changed("partial-payment") {
			body["partial_payment"] = partialPayment
		}
		if currency != "" {
			body["currency"] = currency
		}
		if date > 0 {
			body["date"] = date
		}
		if terms != "" {
			body["terms"] = terms
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
		}

		data, err := client.Patch(basePath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)

	// All states
	updateCmd.Flags().String("description", "", "Invoice description — all states")
	updateCmd.Flags().Int64("expire-by", 0, "Invoice expiry as Unix timestamp — all states")
	updateCmd.Flags().StringArray("note", nil, "Note as key=value (max 15 pairs) — all states")

	// Draft only
	updateCmd.Flags().String("customer-id", "", "Customer ID — draft only")
	updateCmd.Flags().String("customer-name", "", "Customer name — draft only")
	updateCmd.Flags().String("customer-contact", "", "Customer phone number — draft only")
	updateCmd.Flags().String("customer-email", "", "Customer email address — draft only")
	updateCmd.Flags().String("line-items", "", `Line items as a JSON array — draft only`)
	updateCmd.Flags().Bool("sms-notify", true, "Send SMS notification — draft only")
	updateCmd.Flags().Bool("email-notify", true, "Send email notification — draft only")
	updateCmd.Flags().Bool("partial-payment", false, "Allow partial payments — draft only")
	updateCmd.Flags().String("currency", "", "Currency code — draft only")
	updateCmd.Flags().Int64("date", 0, "Invoice issue date as Unix timestamp — draft only")
	updateCmd.Flags().String("terms", "", "Invoice terms — draft only")

	// Billing address (draft only)
	updateCmd.Flags().String("billing-line1", "", "Billing address line 1 — draft only")
	updateCmd.Flags().String("billing-line2", "", "Billing address line 2 — draft only")
	updateCmd.Flags().String("billing-zipcode", "", "Billing address postal code — draft only")
	updateCmd.Flags().String("billing-city", "", "Billing address city — draft only")
	updateCmd.Flags().String("billing-state", "", "Billing address state — draft only")
	updateCmd.Flags().String("billing-country", "", "Billing address country code — draft only")

	// Shipping address (draft only)
	updateCmd.Flags().String("shipping-line1", "", "Shipping address line 1 — draft only")
	updateCmd.Flags().String("shipping-line2", "", "Shipping address line 2 — draft only")
	updateCmd.Flags().String("shipping-zipcode", "", "Shipping address postal code — draft only")
	updateCmd.Flags().String("shipping-city", "", "Shipping address city — draft only")
	updateCmd.Flags().String("shipping-state", "", "Shipping address state — draft only")
	updateCmd.Flags().String("shipping-country", "", "Shipping address country code — draft only")
}
