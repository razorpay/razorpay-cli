package subscriptions

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var offerDeleteCmd = &cobra.Command{
	Use:   "delete-offer <subscription_id> <offer_id>",
	Short: "Delete (unlink) an offer from a subscription",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Delete(basePath + "/" + args[0] + "/" + args[1])
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(offerDeleteCmd)
}
