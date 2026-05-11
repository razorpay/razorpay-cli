package refunds

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <payment_id>",
	Short: "Create a refund for a payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		amount, _ := cmd.Flags().GetInt64("amount")
		speed, _ := cmd.Flags().GetString("speed")
		receipt, _ := cmd.Flags().GetString("receipt")
		notes, _ := cmd.Flags().GetStringArray("note")
		idempotencyKey, _ := cmd.Flags().GetString("idempotency-key")

		body := map[string]interface{}{}
		if amount > 0 {
			body["amount"] = amount
		}
		if speed != "" {
			body["speed"] = speed
		}
		if receipt != "" {
			body["receipt"] = receipt
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		var reqBody any
		if len(body) > 0 {
			reqBody = body
		}

		var data []byte
		var err error
		if idempotencyKey != "" {
			headers := map[string]string{"X-Refund-Idempotency": idempotencyKey}
			data, err = client.PostWithHeaders("/v1/payments/"+args[0]+"/refund", reqBody, headers)
		} else {
			data, err = client.Post("/v1/payments/"+args[0]+"/refund", reqBody)
		}
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(createCmd)

	createCmd.Flags().Int64("amount", 0, "Amount to refund in paise (omit for full refund)")
	createCmd.Flags().String("speed", "", "Refund speed: normal or optimum")
	createCmd.Flags().String("receipt", "", "Unique identifier for internal reference")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	createCmd.Flags().String("idempotency-key", "", "Idempotency key to prevent duplicate refunds (X-Refund-Idempotency)")
}
