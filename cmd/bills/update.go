package bills

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <bill_id>",
	Short: "Update a bill (Billme receipt)",
	Long: `Update a Razorpay Billme receipt.

Scalar fields are passed as typed flags. Nested object and array-of-object
fields (customer, employee, loyalty, line_items, receipt_summary, taxes,
payments, irn) are passed as JSON strings.

Example:
  razorpay bills update bill_xyz \
    --receipt-type tax_invoice \
    --receipt-timestamp 1700000000 \
    --receipt-delivery digital \
    --receipt-summary '{"total":12000,"tax":2160}' \
    --payments '[{"method":"upi","currency":"INR","amount":12000}]'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		receiptType, _ := cmd.Flags().GetString("receipt-type")
		receiptTimestamp, _ := cmd.Flags().GetInt64("receipt-timestamp")
		receiptDelivery, _ := cmd.Flags().GetString("receipt-delivery")

		if receiptTimestamp <= 0 {
			return fmt.Errorf("--receipt-timestamp is required and must be > 0")
		}

		body := map[string]interface{}{
			"receipt_type":      receiptType,
			"receipt_timestamp": receiptTimestamp,
			"receipt_delivery":  receiptDelivery,
		}

		// Optional scalar fields.
		optionalStrings := []struct{ flag, key string }{
			{"store-code", "store_code"},
			{"pos-category", "pos_category"},
		}
		for _, f := range optionalStrings {
			if v, _ := cmd.Flags().GetString(f.flag); v != "" {
				body[f.key] = v
			}
		}

		if tags, _ := cmd.Flags().GetStringArray("tag"); len(tags) > 0 {
			body["tags"] = tags
		}

		// Required nested fields.
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
		data, err := client.Patch(basePath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	// Required scalar flags.
	updateCmd.Flags().String("receipt-type", "", "tax_invoice, sales_invoice, sales_return_invoice, etc.")
	updateCmd.Flags().Int64("receipt-timestamp", 0, "UNIX timestamp when the receipt was generated")
	updateCmd.Flags().String("receipt-delivery", "", "digital, print, digital_and_print")
	_ = updateCmd.MarkFlagRequired("receipt-type")
	_ = updateCmd.MarkFlagRequired("receipt-timestamp")
	_ = updateCmd.MarkFlagRequired("receipt-delivery")

	// Required nested flags (JSON).
	updateCmd.Flags().String("receipt-summary", "", "Receipt summary as JSON object")
	updateCmd.Flags().String("payments", "", "Payments as JSON array of objects with method, currency, amount")
	_ = updateCmd.MarkFlagRequired("receipt-summary")
	_ = updateCmd.MarkFlagRequired("payments")

	// Optional scalar flags.
	updateCmd.Flags().String("store-code", "", "Associated store code")
	updateCmd.Flags().String("pos-category", "", "POS machine type: traditional_pos, kiosk_pos")
	updateCmd.Flags().StringArray("tag", nil, "Tag for the invoice (repeatable)")

	// Optional nested flags (JSON).
	updateCmd.Flags().String("customer", "", "Customer details as JSON object")
	updateCmd.Flags().String("employee", "", "Employee details as JSON array of objects")
	updateCmd.Flags().String("loyalty", "", "Customer loyalty details as JSON object")
	updateCmd.Flags().String("line-items", "", "Line items as JSON array of objects")
	updateCmd.Flags().String("taxes", "", "Taxes as JSON array of objects")
	updateCmd.Flags().String("irn", "", "Invoice Reference Number details as JSON object")

	Cmd.AddCommand(updateCmd)
}
