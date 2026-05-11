package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferUpdateCmd = &cobra.Command{
	Use:   "update <transfer_id>",
	Short: "Modify settlement hold for a transfer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("on-hold") {
			return fmt.Errorf("--on-hold is required")
		}
		client := cmdutil.NewClient()
		onHold, _ := cmd.Flags().GetBool("on-hold")
		body := map[string]any{
			"on_hold": onHold,
		}
		if onHoldUntil, _ := cmd.Flags().GetInt64("on-hold-until"); onHoldUntil > 0 {
			body["on_hold_until"] = onHoldUntil
		}

		data, err := client.Patch(transfersPath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferUpdateCmd)

	transferUpdateCmd.Flags().Bool("on-hold", false, "Whether to hold (true) or release (false) the settlement")
	transferUpdateCmd.Flags().Int64("on-hold-until", 0, "Unix timestamp until which settlement is held (optional)")
}
