package route

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFromPaymentCmd = &cobra.Command{
	Use:   "create-from-payment <payment_id>",
	Short: "Create transfers from a captured payment",
	Long: `Create one or more transfers from a captured payment to linked accounts.

	Pass the transfers as a JSON array via --transfers. Each element supports:
  	account, amount, currency, notes, linked_account_notes, on_hold, on_hold_until

	Example:
  	razorpay route transfers create-from-payment pay_ABC123 \
    --transfers '[
      {"account":"acc_ABC","amount":10000,"currency":"INR","notes":{"branch":"south"},"linked_account_notes":["branch"]},
      {"account":"acc_XYZ","amount":20000,"currency":"INR","on_hold":true}
    ]'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		transfersJSON, _ := cmd.Flags().GetString("transfers")

		if transfersJSON == "" {
			return fmt.Errorf("--transfers is required (JSON array of transfer objects)")
		}

		var transfers any
		if err := json.Unmarshal([]byte(transfersJSON), &transfers); err != nil {
			return fmt.Errorf("--transfers is not valid JSON: %w", err)
		}

		body := map[string]any{
			"transfers": transfers,
		}

		data, err := client.Post("/v1/payments/"+args[0]+"/transfers", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFromPaymentCmd)
	transferFromPaymentCmd.Flags().String("transfers", "", `Transfers as a JSON array (required). Each object supports: account, amount, currency, notes, linked_account_notes, on_hold, on_hold_until`)
}
