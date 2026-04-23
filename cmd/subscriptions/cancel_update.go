package subscriptions

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var cancelUpdateCmd = &cobra.Command{
	Use:   "cancel-update <subscription_id>",
	Short: "Cancel pending scheduled changes for a subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Post(basePath+"/"+args[0]+"/cancel_scheduled_changes", nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(cancelUpdateCmd)
}
