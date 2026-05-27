package invoices

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <invoice_id>",
	Short: "Delete a draft invoice",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Delete(basePath + "/" + args[0])
		if err != nil {
			cmdutil.HandleErr(err)
		}
		if len(data) > 0 {
			api.PrettyPrint(data)
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(deleteCmd)
}
