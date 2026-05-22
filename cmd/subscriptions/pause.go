package subscriptions

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause <subscription_id>",
	Short: "Pause a subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pauseAt, _ := cmd.Flags().GetString("pause-at")
		if pauseAt != "now" && pauseAt != "cycle_end" {
			return fmt.Errorf("--pause-at must be now or cycle_end")
		}
		client := cmdutil.NewClient()
		body := map[string]interface{}{
			"pause_at": pauseAt,
		}
		data, err := client.Post(basePath+"/"+args[0]+"/pause", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(pauseCmd)

	pauseCmd.Flags().String("pause-at", "now", "When to pause: now")
}
