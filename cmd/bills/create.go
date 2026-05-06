package bills

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a bill (Billme receipt)",
	Long: `Create a Razorpay Billme receipt.

Scalar fields are passed as typed flags. Nested object and array-of-object
fields (customer, employee, loyalty, line_items, receipt_summary, taxes,
payments, event, irn) are passed as JSON strings.

Example:
  razorpay bills create \
    --business-type retail \
    --business-category food_and_beverages \
    --receipt-timestamp 1700000000 \
    --receipt-number INV-001 \
    --receipt-type tax_invoice \
    --receipt-delivery digital \
    --receipt-summary '{"total":10000,"tax":1800,"payment_status":"paid"}' \
    --payments '[{"method":"card","currency":"INR","amount":10000}]' \
    --customer '{"name":"Alice","email":"alice@example.com"}' \
    --line-items '[{"sku":"X","quantity":1,"price":10000}]'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		businessType, _ := cmd.Flags().GetString("business-type")
		businessCategory, _ := cmd.Flags().GetString("business-category")
		receiptTimestamp, _ := cmd.Flags().GetInt64("receipt-timestamp")
		receiptNumber, _ := cmd.Flags().GetString("receipt-number")
		receiptType, _ := cmd.Flags().GetString("receipt-type")
		receiptDelivery, _ := cmd.Flags().GetString("receipt-delivery")

		if receiptTimestamp <= 0 {
			return fmt.Errorf("--receipt-timestamp is required and must be > 0")
		}

		body := map[string]interface{}{
			"business_type":     businessType,
			"business_category": businessCategory,
			"receipt_timestamp": receiptTimestamp,
			"receipt_number":    receiptNumber,
			"receipt_type":      receiptType,
			"receipt_delivery":  receiptDelivery,
		}

		// Optional scalar fields.
		optionalStrings := []struct{ flag, key string }{
			{"store-code", "store_code"},
			{"bar-code-number", "bar_code_number"},
			{"qr-code-number", "qr_code_number"},
			{"billing-pos-number", "billing_pos_number"},
			{"pos-category", "pos_category"},
			{"order-number", "order_number"},
			{"order-service-type", "order_service_type"},
			{"delivery-status-url", "delivery_status_url"},
		}
		for _, f := range optionalStrings {
			if v, _ := cmd.Flags().GetString(f.flag); v != "" {
				body[f.key] = v
			}
		}

		if tags, _ := cmd.Flags().GetStringArray("tag"); len(tags) > 0 {
			body["tags"] = tags
		}

		// Required nested fields (must be provided).
		for _, f := range []struct{ flag, key string }{
			{"receipt-summary", "receipt_summary"},
			{"payments", "payments"},
		} {
			v, err := parseJSONFlag(cmd, f.flag)
			if err != nil {
				return err
			}
			if v == nil {
				return fmt.Errorf("--%s is required (JSON)", f.flag)
			}
			body[f.key] = v
		}

		// Optional nested fields.
		for _, f := range []struct{ flag, key string }{
			{"customer", "customer"},
			{"employee", "employee"},
			{"loyalty", "loyalty"},
			{"line-items", "line_items"},
			{"taxes", "taxes"},
			{"event", "event"},
			{"irn", "irn"},
		} {
			v, err := parseJSONFlag(cmd, f.flag)
			if err != nil {
				return err
			}
			if v != nil {
				body[f.key] = v
			}
		}

		client := cmdutil.NewClient()
		data, err := client.Post(basePath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	// Required scalar flags.
	createCmd.Flags().String("business-type", "", "Type of business: ecommerce, retail")
	createCmd.Flags().String("business-category", "", "Business category: events, food_and_beverages, etc.")
	createCmd.Flags().Int64("receipt-timestamp", 0, "UNIX timestamp when the receipt was generated")
	createCmd.Flags().String("receipt-number", "", "Unique receipt number for the bill")
	createCmd.Flags().String("receipt-type", "", "tax_invoice, sales_invoice, sales_return_invoice, etc.")
	createCmd.Flags().String("receipt-delivery", "", "digital, print, digital_and_print")
	_ = createCmd.MarkFlagRequired("business-type")
	_ = createCmd.MarkFlagRequired("business-category")
	_ = createCmd.MarkFlagRequired("receipt-timestamp")
	_ = createCmd.MarkFlagRequired("receipt-number")
	_ = createCmd.MarkFlagRequired("receipt-type")
	_ = createCmd.MarkFlagRequired("receipt-delivery")

	// Required nested flags (JSON).
	createCmd.Flags().String("receipt-summary", "", "Receipt summary as JSON object (totals, taxes, discounts, payment_status)")
	createCmd.Flags().String("payments", "", "Payments as JSON array of objects with method, currency, amount")
	_ = createCmd.MarkFlagRequired("receipt-summary")
	_ = createCmd.MarkFlagRequired("payments")

	// Optional scalar flags.
	createCmd.Flags().String("store-code", "", "Associated store code (required if multi-store setup)")
	createCmd.Flags().String("bar-code-number", "", "Bar code generated after the transaction")
	createCmd.Flags().String("qr-code-number", "", "QR code generated after the transaction")
	createCmd.Flags().String("billing-pos-number", "", "POS number of the machine that generated the bill")
	createCmd.Flags().String("pos-category", "", "POS machine type: traditional_pos, kiosk_pos")
	createCmd.Flags().String("order-number", "", "Incremental order number of the generated bill")
	createCmd.Flags().String("order-service-type", "", "Order service type: dine_in, take_away")
	createCmd.Flags().String("delivery-status-url", "", "Order delivery status URL (ecommerce)")
	createCmd.Flags().StringArray("tag", nil, "Tag for the invoice (repeatable)")

	// Optional nested flags (JSON).
	createCmd.Flags().String("customer", "", "Customer details as JSON object (required if receipt-delivery is digital)")
	createCmd.Flags().String("employee", "", "Employee details as JSON array of objects")
	createCmd.Flags().String("loyalty", "", "Customer loyalty details as JSON object")
	createCmd.Flags().String("line-items", "", "Line items as JSON array of objects (required for some receipt types)")
	createCmd.Flags().String("taxes", "", "Taxes as JSON array of objects (required for tax_invoice and similar)")
	createCmd.Flags().String("event", "", "Event booking details as JSON object (required if business-category is events)")
	createCmd.Flags().String("irn", "", "Invoice Reference Number details as JSON object")

	Cmd.AddCommand(createCmd)
}
