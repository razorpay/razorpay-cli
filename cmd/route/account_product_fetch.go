package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var productFetchCmd = &cobra.Command{
	Use:   "product-fetch <account_id> <product_id>",
	Short: "Fetch Route product configuration for a linked account",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get(accountsPath+"/"+args[0]+"/products/"+args[1], nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(productFetchCmd)
}
