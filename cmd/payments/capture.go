package payments

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var captureCmd = &cobra.Command{
	Use:   "capture <payment_id>",
	Short: "Capture an authorized payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		if amount <= 0 {
			return fmt.Errorf("--amount is required and must be > 0")
		}
		body := map[string]interface{}{
			"amount":   amount,
			"currency": currency,
		}
		data, err := client.Post(basePath+"/"+args[0]+"/capture", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(captureCmd)

	captureCmd.Flags().Int("amount", 0, "Amount to capture in smallest currency unit (e.g. paise)")
	captureCmd.Flags().String("currency", "INR", "Currency code (e.g. INR)")
	_ = captureCmd.MarkFlagRequired("amount")
}
