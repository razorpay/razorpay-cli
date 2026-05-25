package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFetchCmd = &cobra.Command{
	Use:   "fetch <transfer_id>",
	Short: "Fetch a transfer by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get(transfersPath+"/"+args[0], nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFetchCmd)
}
