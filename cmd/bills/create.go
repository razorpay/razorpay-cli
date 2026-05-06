package bills

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new bill",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		// Required scalar flags
		businessType, _ := cmd.Flags().GetString("business-type")
		businessCategory, _ := cmd.Flags().GetString("business-category")
		receiptTimestamp, _ := cmd.Flags().GetInt64("receipt-timestamp")
		receiptNumber, _ := cmd.Flags().GetString("receipt-number")
		receiptType, _ := cmd.Flags().GetString("receipt-type")
		receiptDelivery, _ := cmd.Flags().GetString("receipt-delivery")

		// Required JSON object/array flags
		receiptSummaryJSON, _ := cmd.Flags().GetString("receipt-summary")
		paymentsJSON, _ := cmd.Flags().GetString("payments")

		// Optional scalar flags
		storeCode, _ := cmd.Flags().GetString("store-code")
		barCodeNumber, _ := cmd.Flags().GetString("bar-code-number")
		qrCodeNumber, _ := cmd.Flags().GetString("qr-code-number")
		billingPosNumber, _ := cmd.Flags().GetString("billing-pos-number")
		posCategory, _ := cmd.Flags().GetString("pos-category")
		orderNumber, _ := cmd.Flags().GetString("order-number")
		orderServiceType, _ := cmd.Flags().GetString("order-service-type")
		deliveryStatusURL, _ := cmd.Flags().GetString("delivery-status-url")
		tags, _ := cmd.Flags().GetStringArray("tag")

		// Optional JSON object/array flags
		customerJSON, _ := cmd.Flags().GetString("customer")
		employeeJSON, _ := cmd.Flags().GetString("employee")
		loyaltyJSON, _ := cmd.Flags().GetString("loyalty")
		lineItemsJSON, _ := cmd.Flags().GetString("line-items")
		taxesJSON, _ := cmd.Flags().GetString("taxes")
		eventJSON, _ := cmd.Flags().GetString("event")
		irnJSON, _ := cmd.Flags().GetString("irn")

		// Validate required scalars
		if businessType == "" {
			return fmt.Errorf("--business-type is required (ecommerce or retail)")
		}
		if businessCategory == "" {
			return fmt.Errorf("--business-category is required (events, food_and_beverages, retail_and_consumer_goods, or other_services)")
		}
		if receiptTimestamp <= 0 {
			return fmt.Errorf("--receipt-timestamp is required and must be a positive Unix timestamp")
		}
		if receiptNumber == "" {
			return fmt.Errorf("--receipt-number is required")
		}
		if receiptType == "" {
			return fmt.Errorf("--receipt-type is required")
		}
		if receiptDelivery == "" {
			return fmt.Errorf("--receipt-delivery is required (digital, print, or digital_and_print)")
		}

		body := map[string]any{
			"business_type":     businessType,
			"business_category": businessCategory,
			"receipt_timestamp": receiptTimestamp,
			"receipt_number":    receiptNumber,
			"receipt_type":      receiptType,
			"receipt_delivery":  receiptDelivery,
		}

		// Required JSON: receipt_summary
		if receiptSummaryJSON == "" {
			return fmt.Errorf("--receipt-summary is required (JSON object)")
		}
		var receiptSummary any
		if err := json.Unmarshal([]byte(receiptSummaryJSON), &receiptSummary); err != nil {
			return fmt.Errorf("--receipt-summary is not valid JSON: %w", err)
		}
		body["receipt_summary"] = receiptSummary

		// Required JSON: payments
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
		if barCodeNumber != "" {
			body["bar_code_number"] = barCodeNumber
		}
		if qrCodeNumber != "" {
			body["qr_code_number"] = qrCodeNumber
		}
		if billingPosNumber != "" {
			body["billing_pos_number"] = billingPosNumber
		}
		if posCategory != "" {
			body["pos_category"] = posCategory
		}
		if orderNumber != "" {
			body["order_number"] = orderNumber
		}
		if orderServiceType != "" {
			body["order_service_type"] = orderServiceType
		}
		if deliveryStatusURL != "" {
			body["delivery_status_url"] = deliveryStatusURL
		}
		if len(tags) > 0 {
			body["tags"] = tags
		}

		// Optional JSON objects/arrays
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
		if eventJSON != "" {
			var event any
			if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
				return fmt.Errorf("--event is not valid JSON: %w", err)
			}
			body["event"] = event
		}
		if irnJSON != "" {
			var irn any
			if err := json.Unmarshal([]byte(irnJSON), &irn); err != nil {
				return fmt.Errorf("--irn is not valid JSON: %w", err)
			}
			body["irn"] = irn
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

	// Required scalars
	createCmd.Flags().String("business-type", "", "Business type: ecommerce or retail (required)")
	createCmd.Flags().String("business-category", "", "Business category: events, food_and_beverages, retail_and_consumer_goods, other_services (required)")
	createCmd.Flags().Int64("receipt-timestamp", 0, "Unix timestamp of when the receipt was generated (required)")
	createCmd.Flags().String("receipt-number", "", "Unique receipt number for the bill (required)")
	createCmd.Flags().String("receipt-type", "", "Receipt type: tax_invoice, sales_invoice, sales_return_invoice, proforma_invoice, credit_invoice, purchase_invoice, debit_invoice, order_confirmation (required)")
	createCmd.Flags().String("receipt-delivery", "", "Receipt delivery: digital, print, or digital_and_print (required)")

	// Required JSON
	createCmd.Flags().String("receipt-summary", "", `Receipt summary as JSON object (required, e.g. '{"total_quantity":1,"sub_total_amount":100,"currency":"INR","net_payable_amount":118}')`)
	createCmd.Flags().String("payments", "", `Payments as JSON array (required, e.g. '[{"method":"card","currency":"INR","amount":118}]')`)

	_ = createCmd.MarkFlagRequired("business-type")
	_ = createCmd.MarkFlagRequired("business-category")
	_ = createCmd.MarkFlagRequired("receipt-timestamp")
	_ = createCmd.MarkFlagRequired("receipt-number")
	_ = createCmd.MarkFlagRequired("receipt-type")
	_ = createCmd.MarkFlagRequired("receipt-delivery")
	_ = createCmd.MarkFlagRequired("receipt-summary")
	_ = createCmd.MarkFlagRequired("payments")

	// Optional scalars
	createCmd.Flags().String("store-code", "", "Associated store code (required for multi-store setup)")
	createCmd.Flags().String("bar-code-number", "", "Bar code generated after the transaction (digital bills)")
	createCmd.Flags().String("qr-code-number", "", "QR code generated after the transaction (digital bills)")
	createCmd.Flags().String("billing-pos-number", "", "POS number of the machine that generated the bill (retail)")
	createCmd.Flags().String("pos-category", "", "POS machine type: traditional_pos or kiosk_pos (retail)")
	createCmd.Flags().String("order-number", "", "Incremental order number of the bill")
	createCmd.Flags().String("order-service-type", "", "Order service type: dine_in or take_away (food_and_beverages)")
	createCmd.Flags().String("delivery-status-url", "", "Order delivery status URL (ecommerce)")
	createCmd.Flags().StringArray("tag", nil, "Tag associated with the invoice (repeatable)")

	// Optional JSON
	createCmd.Flags().String("customer", "", "Customer details as JSON object (contact, name, email, customer_id, billing_address, shipping_address, ...)")
	createCmd.Flags().String("employee", "", "Employees as JSON array of objects (id, name, role)")
	createCmd.Flags().String("loyalty", "", "Customer loyalty details as JSON object")
	createCmd.Flags().String("line-items", "", `Line items as JSON array (required if receipt_type is not credit_invoice or debit_invoice)`)
	createCmd.Flags().String("taxes", "", `Taxes as JSON array (required for tax_invoice, purchase_invoice, sales_invoice)`)
	createCmd.Flags().String("event", "", "Event booking details as JSON object (required if business_category is events)")
	createCmd.Flags().String("irn", "", "IRN details as JSON object (acknowledgement_number, acknowledgement_date, qr_code, irn_number)")
}
