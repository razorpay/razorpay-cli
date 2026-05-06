package bills

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <bill_id>",
	Short: "Update the details of a bill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		// Required scalar flags
		receiptType, _ := cmd.Flags().GetString("receipt-type")
		receiptTimestamp, _ := cmd.Flags().GetInt64("receipt-timestamp")
		receiptDelivery, _ := cmd.Flags().GetString("receipt-delivery")

		// Required JSON
		receiptSummaryJSON, _ := cmd.Flags().GetString("receipt-summary")
		paymentsJSON, _ := cmd.Flags().GetString("payments")

		// Optional scalar flags
		storeCode, _ := cmd.Flags().GetString("store-code")
		posCategory, _ := cmd.Flags().GetString("pos-category")
		tags, _ := cmd.Flags().GetStringArray("tag")

		// Optional JSON flags
		customerJSON, _ := cmd.Flags().GetString("customer")
		employeeJSON, _ := cmd.Flags().GetString("employee")
		loyaltyJSON, _ := cmd.Flags().GetString("loyalty")
		lineItemsJSON, _ := cmd.Flags().GetString("line-items")
		taxesJSON, _ := cmd.Flags().GetString("taxes")
		irnJSON, _ := cmd.Flags().GetString("irn")

		// Validate required scalars
		if receiptType == "" {
			return fmt.Errorf("--receipt-type is required")
		}
		if receiptTimestamp <= 0 {
			return fmt.Errorf("--receipt-timestamp is required and must be a positive Unix timestamp")
		}
		if receiptDelivery == "" {
			return fmt.Errorf("--receipt-delivery is required (digital, print, or digital_and_print)")
		}

		body := map[string]any{
			"receipt_type":      receiptType,
			"receipt_timestamp": receiptTimestamp,
			"receipt_delivery":  receiptDelivery,
		}

		// Required JSON
		if receiptSummaryJSON == "" {
			return fmt.Errorf("--receipt-summary is required (JSON object)")
		}
		var receiptSummary any
		if err := json.Unmarshal([]byte(receiptSummaryJSON), &receiptSummary); err != nil {
			return fmt.Errorf("--receipt-summary is not valid JSON: %w", err)
		}
		body["receipt_summary"] = receiptSummary

		if paymentsJSON == "" {
			return fmt.Errorf("--payments is required (JSON array)")
		}
		var payments any
		if err := json.Unmarshal([]byte(paymentsJSON), &payments); err != nil {
			return fmt.Errorf("--payments is not valid JSON: %w", err)
		}
		body["payments"] = payments

		// Optional scalars
		if storeCode != "" {
			body["store_code"] = storeCode
		}
		if posCategory != "" {
			body["pos_category"] = posCategory
		}
		if len(tags) > 0 {
			body["tags"] = tags
		}

		// Optional JSON
		if customerJSON != "" {
			var customer any
			if err := json.Unmarshal([]byte(customerJSON), &customer); err != nil {
				return fmt.Errorf("--customer is not valid JSON: %w", err)
			}
			body["customer"] = customer
		}
		if employeeJSON != "" {
			var employee any
			if err := json.Unmarshal([]byte(employeeJSON), &employee); err != nil {
				return fmt.Errorf("--employee is not valid JSON: %w", err)
			}
			body["employee"] = employee
		}
		if loyaltyJSON != "" {
			var loyalty any
			if err := json.Unmarshal([]byte(loyaltyJSON), &loyalty); err != nil {
				return fmt.Errorf("--loyalty is not valid JSON: %w", err)
			}
			body["loyalty"] = loyalty
		}
		if lineItemsJSON != "" {
			var lineItems any
			if err := json.Unmarshal([]byte(lineItemsJSON), &lineItems); err != nil {
				return fmt.Errorf("--line-items is not valid JSON: %w", err)
			}
			body["line_items"] = lineItems
		}
		if taxesJSON != "" {
			var taxes any
			if err := json.Unmarshal([]byte(taxesJSON), &taxes); err != nil {
				return fmt.Errorf("--taxes is not valid JSON: %w", err)
			}
			body["taxes"] = taxes
		}
		if irnJSON != "" {
			var irn any
			if err := json.Unmarshal([]byte(irnJSON), &irn); err != nil {
				return fmt.Errorf("--irn is not valid JSON: %w", err)
			}
			body["irn"] = irn
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

	updateCmd.Flags().String("receipt-type", "", "Receipt type: tax_invoice, sales_invoice, sales_return_invoice, proforma_invoice, credit_invoice, purchase_invoice, debit_invoice (required)")
	updateCmd.Flags().Int64("receipt-timestamp", 0, "Unix timestamp of when the receipt was generated (required)")
	updateCmd.Flags().String("receipt-delivery", "", "Receipt delivery: digital, print, or digital_and_print (required)")
	updateCmd.Flags().String("receipt-summary", "", "Receipt summary as JSON object (required)")
	updateCmd.Flags().String("payments", "", "Payments as JSON array (required)")

	_ = updateCmd.MarkFlagRequired("receipt-type")
	_ = updateCmd.MarkFlagRequired("receipt-timestamp")
	_ = updateCmd.MarkFlagRequired("receipt-delivery")
	_ = updateCmd.MarkFlagRequired("receipt-summary")
	_ = updateCmd.MarkFlagRequired("payments")

	updateCmd.Flags().String("store-code", "", "Associated store code")
	updateCmd.Flags().String("pos-category", "", "POS machine type: traditional_pos or kiosk_pos (retail)")
	updateCmd.Flags().StringArray("tag", nil, "Tag associated with the invoice (repeatable)")

	updateCmd.Flags().String("customer", "", "Customer details as JSON object")
	updateCmd.Flags().String("employee", "", "Employees as JSON array of objects (id, name, role)")
	updateCmd.Flags().String("loyalty", "", "Customer loyalty details as JSON object")
	updateCmd.Flags().String("line-items", "", "Line items as JSON array")
	updateCmd.Flags().String("taxes", "", "Taxes as JSON array (name, percentage, amount)")
	updateCmd.Flags().String("irn", "", "IRN details as JSON object")
}
