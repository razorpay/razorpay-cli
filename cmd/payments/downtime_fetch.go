package payments

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var downtimeFetchCmd = &cobra.Command{
	Use:   "fetch <downtime_id>",
	Short: "Fetch a payment downtime by ID",
	Long: `Fetch details of a specific payment method downtime by its ID.

Example:
  razorpay payments downtime fetch down_F7LroRQAAFuswd`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return fmt.Errorf("downtime_id is required")
		}
		client := cmdutil.NewClient()
		data, err := client.Get(basePath+"/downtimes/"+args[0], nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	downtimeCmd.AddCommand(downtimeFetchCmd)
}
