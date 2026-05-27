package subscriptions

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel <subscription_id>",
	Short: "Cancel a subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		cancelAtCycleEnd, _ := cmd.Flags().GetBool("cancel-at-cycle-end")
		body := map[string]interface{}{
			"cancel_at_cycle_end": cancelAtCycleEnd,
		}
		data, err := client.Post(basePath+"/"+args[0]+"/cancel", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(cancelCmd)

	cancelCmd.Flags().Bool("cancel-at-cycle-end", false, "Cancel at end of current billing cycle instead of immediately")
}
