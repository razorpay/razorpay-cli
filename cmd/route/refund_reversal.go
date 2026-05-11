package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var refundReversalCmd = &cobra.Command{
	Use:   "refund-with-reversal <payment_id>",
	Short: "Refund a payment and reverse all linked transfers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		amount, _ := cmd.Flags().GetInt64("amount")
		if amount == 0 {
			return fmt.Errorf("--amount is required")
		}

		client := cmdutil.NewClient()
		body := map[string]any{
			"amount":      amount,
			"reverse_all": 1,
		}
		if notes, _ := cmd.Flags().GetStringArray("note"); len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		data, err := client.Post("/v1/payments/"+args[0]+"/refund", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(refundReversalCmd)

	refundReversalCmd.Flags().Int64("amount", 0, "Refund amount in paise (required)")
	refundReversalCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
