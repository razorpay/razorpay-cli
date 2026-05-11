package payments

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var downtimeListCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch all payment method downtimes",
	Long: `Fetch all ongoing and scheduled payment method downtimes.

Example:
  razorpay payments downtime list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get(basePath+"/downtimes", nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	downtimeCmd.AddCommand(downtimeListCmd)
}
