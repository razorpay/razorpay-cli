package smartcollect

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var tpvAddPayerCmd = &cobra.Command{
	Use:   "tpv-add-payer <virtual_account_id>",
	Short: "Add an allowed payer bank account to a TPV virtual account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ifsc, _ := cmd.Flags().GetString("ifsc")
		accountNumber, _ := cmd.Flags().GetString("account-number")
		client := cmdutil.NewClient()

		body := map[string]any{
			"type": "bank_account",
			"bank_account": map[string]any{
				"ifsc":           ifsc,
				"account_number": accountNumber,
			},
		}

		data, err := client.Post(basePath+"/"+args[0]+"/allowed_payers", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(tpvAddPayerCmd)

	tpvAddPayerCmd.Flags().String("ifsc", "", "Bank IFSC code of the allowed payer (required)")
	tpvAddPayerCmd.Flags().String("account-number", "", "Bank account number of the allowed payer (required)")
	_ = tpvAddPayerCmd.MarkFlagRequired("ifsc")
	_ = tpvAddPayerCmd.MarkFlagRequired("account-number")
}
