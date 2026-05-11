package orders

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new order",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		receipt, _ := cmd.Flags().GetString("receipt")
		notes, _ := cmd.Flags().GetStringArray("note")

		if amount <= 0 {
			return fmt.Errorf("--amount is required and must be > 0")
		}
		if currency == "" {
			return fmt.Errorf("--currency is required")
		}

		body := map[string]interface{}{
			"amount":   amount,
			"currency": currency,
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

	createCmd.Flags().Int("amount", 0, "Order amount in smallest currency unit (e.g. paise for INR)")
	createCmd.Flags().String("currency", "INR", "Currency code (e.g. INR)")
	createCmd.Flags().String("receipt", "", "Receipt number for your internal reference (max 40 chars)")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	_ = createCmd.MarkFlagRequired("amount")
}
