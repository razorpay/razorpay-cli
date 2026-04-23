package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var accountFetchCmd = &cobra.Command{
	Use:   "fetch <account_id>",
	Short: "Fetch a linked account by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get(accountsPath+"/"+args[0], nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(accountFetchCmd)
}
