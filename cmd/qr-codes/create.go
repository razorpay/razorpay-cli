package qrcodes

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		qrType, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		usage, _ := cmd.Flags().GetString("usage")
		fixedAmount, _ := cmd.Flags().GetBool("fixed-amount")
		paymentAmount, _ := cmd.Flags().GetInt64("payment-amount")
		description, _ := cmd.Flags().GetString("description")
		customerID, _ := cmd.Flags().GetString("customer-id")
		closeBy, _ := cmd.Flags().GetInt64("close-by")
		notes, _ := cmd.Flags().GetStringArray("note")

		if qrType == "" {
			return fmt.Errorf("--type is required (upi_qr)")
		}
		if usage == "" {
			return fmt.Errorf("--usage is required (single_use or multiple_use)")
		}
		if fixedAmount && paymentAmount <= 0 {
			return fmt.Errorf("--payment-amount is required when --fixed-amount is set")
		}

		body := map[string]any{
			"type":  qrType,
			"usage": usage,
		}
		if name != "" {
			body["name"] = name
		}
		if cmd.Flags().Changed("fixed-amount") {
			body["fixed_amount"] = fixedAmount
		}
		if paymentAmount > 0 {
			body["payment_amount"] = paymentAmount
		}
		if description != "" {
			body["description"] = description
		}
		if customerID != "" {
			body["customer_id"] = customerID
		}
		if closeBy > 0 {
			body["close_by"] = closeBy
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

	createCmd.Flags().String("type", "", "QR code type: upi_qr (required)")
	createCmd.Flags().String("name", "", "Label for the QR code (e.g. Store Front Display)")
	createCmd.Flags().String("usage", "", "Payment acceptance: single_use or multiple_use (required)")
	createCmd.Flags().Bool("fixed-amount", false, "Accept only a specific fixed amount")
	createCmd.Flags().Int64("payment-amount", 0, "Fixed amount in paise (required with --fixed-amount)")
	createCmd.Flags().String("description", "", "Brief description of the QR code's purpose")
	createCmd.Flags().String("customer-id", "", "Customer ID to link to this QR code")
	createCmd.Flags().Int64("close-by", 0, "Unix timestamp for automatic closure (single_use only, 2–120 min from now)")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	_ = createCmd.MarkFlagRequired("type")
}
