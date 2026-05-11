package invoices

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new invoice",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		invoiceType, _ := cmd.Flags().GetString("type")
		customerID, _ := cmd.Flags().GetString("customer-id")
		description, _ := cmd.Flags().GetString("description")
		draft, _ := cmd.Flags().GetBool("draft")
		expireBy, _ := cmd.Flags().GetInt64("expire-by")
		date, _ := cmd.Flags().GetInt64("date")
		currency, _ := cmd.Flags().GetString("currency")
		terms, _ := cmd.Flags().GetString("terms")
		smsNotify, _ := cmd.Flags().GetBool("sms-notify")
		emailNotify, _ := cmd.Flags().GetBool("email-notify")
		partialPayment, _ := cmd.Flags().GetBool("partial-payment")
		lineItemsJSON, _ := cmd.Flags().GetString("line-items")
		notes, _ := cmd.Flags().GetStringArray("note")
		customerName, _ := cmd.Flags().GetString("customer-name")
		customerContact, _ := cmd.Flags().GetString("customer-contact")
		customerEmail, _ := cmd.Flags().GetString("customer-email")

		if invoiceType == "" {
			return fmt.Errorf("--type is required (invoice, link, or ecod)")
		}

		body := map[string]any{
			"type": invoiceType,
		}
		if customerID != "" {
			body["customer_id"] = customerID
		}
		if description != "" {
			body["description"] = description
		}
		if draft {
			body["draft"] = "1"
		}
		if expireBy > 0 {
			body["expire_by"] = expireBy
		}
		if date > 0 {
			body["date"] = date
		}
		if currency != "" {
			body["currency"] = currency
		}
		if terms != "" {
			body["terms"] = terms
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

		if lineItemsJSON != "" {
			var items any
			if err := json.Unmarshal([]byte(lineItemsJSON), &items); err != nil {
				return fmt.Errorf("--line-items is not valid JSON: %w", err)
			}
			body["line_items"] = items
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

		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
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

	createCmd.Flags().String("type", "", "Invoice type: invoice, link, or ecod (required)")
	createCmd.Flags().String("description", "", "Invoice description (max 2048 characters)")
	createCmd.Flags().Bool("draft", false, "Create invoice in draft state")
	createCmd.Flags().Int64("expire-by", 0, "Invoice expiry as Unix timestamp")
	createCmd.Flags().Int64("date", 0, "Invoice issue date as Unix timestamp")
	createCmd.Flags().String("currency", "", "Currency code (required for international payments)")
	createCmd.Flags().String("terms", "", "Invoice terms and conditions (max 2048 characters)")
	createCmd.Flags().Bool("sms-notify", true, "Send SMS notification to customer")
	createCmd.Flags().Bool("email-notify", true, "Send email notification to customer")
	createCmd.Flags().Bool("partial-payment", false, "Allow customer to make partial payments")
	createCmd.Flags().String("line-items", "", `Line items as a JSON array (e.g. '[{"name":"Service","amount":5000,"currency":"INR"}]')`)
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")

	// Customer — use --customer-id for existing customers, or inline fields for new ones
	createCmd.Flags().String("customer-id", "", "Existing customer ID to link to the invoice")
	createCmd.Flags().String("customer-name", "", "Customer name")
	createCmd.Flags().String("customer-contact", "", "Customer phone number")
	createCmd.Flags().String("customer-email", "", "Customer email address")

	// Billing address
	createCmd.Flags().String("billing-line1", "", "Billing address line 1")
	createCmd.Flags().String("billing-line2", "", "Billing address line 2")
	createCmd.Flags().String("billing-zipcode", "", "Billing address postal code")
	createCmd.Flags().String("billing-city", "", "Billing address city")
	createCmd.Flags().String("billing-state", "", "Billing address state")
	createCmd.Flags().String("billing-country", "", "Billing address country code (e.g. in)")

	// Shipping address
	createCmd.Flags().String("shipping-line1", "", "Shipping address line 1")
	createCmd.Flags().String("shipping-line2", "", "Shipping address line 2")
	createCmd.Flags().String("shipping-zipcode", "", "Shipping address postal code")
	createCmd.Flags().String("shipping-city", "", "Shipping address city")
	createCmd.Flags().String("shipping-state", "", "Shipping address state")
	createCmd.Flags().String("shipping-country", "", "Shipping address country code (e.g. in)")
}
