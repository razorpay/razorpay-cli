package smartcollect

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var tpvDeletePayerCmd = &cobra.Command{
	Use:   "tpv-delete-payer <virtual_account_id> <allowed_payer_id>",
	Short: "Delete an allowed payer from a TPV virtual account",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Delete(basePath + "/" + args[0] + "/allowed_payers/" + args[1])
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(tpvDeletePayerCmd)
}
