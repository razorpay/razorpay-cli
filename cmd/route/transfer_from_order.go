package route

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFromOrderCmd = &cobra.Command{
	Use:   "create-from-order",
	Short: "Create an order with embedded transfers to linked accounts",
	Long: `Create an order with one or more transfers to linked accounts.

	Pass the transfers as a JSON array via --transfers. Each element supports:
  	account, amount, currency, notes, linked_account_notes, on_hold, on_hold_until

	Example:
  	razorpay route transfers create-from-order --amount 50000 \
    --transfers '[
      {"account":"acc_ABC","amount":10000,"currency":"INR","notes":{"branch":"south"},"linked_account_notes":["branch"]},
      {"account":"acc_XYZ","amount":20000,"currency":"INR","on_hold":true}
    ]'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		receipt, _ := cmd.Flags().GetString("receipt")
		transfersJSON, _ := cmd.Flags().GetString("transfers")

		if amount == 0 {
			return fmt.Errorf("--amount is required")
		}
		if transfersJSON == "" {
			return fmt.Errorf("--transfers is required (JSON array of transfer objects)")
		}

		var transfers any
		if err := json.Unmarshal([]byte(transfersJSON), &transfers); err != nil {
			return fmt.Errorf("--transfers is not valid JSON: %w", err)
		}

		body := map[string]any{
			"amount":    amount,
			"currency":  currency,
			"transfers": transfers,
		}
		if receipt != "" {
			body["receipt"] = receipt
		}

		data, err := client.Post("/v1/orders", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFromOrderCmd)

	transferFromOrderCmd.Flags().Int64("amount", 0, "Order amount in paise (required)")
	transferFromOrderCmd.Flags().String("currency", "INR", "Order currency code (default INR)")
	transferFromOrderCmd.Flags().String("receipt", "", "Order receipt / reference identifier")
	transferFromOrderCmd.Flags().String("transfers", "", `Transfers as a JSON array (required). Each object supports: account, amount, currency, notes, linked_account_notes, on_hold, on_hold_until`)
}
