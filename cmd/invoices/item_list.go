package invoices

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var itemListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		if from, _ := cmd.Flags().GetInt64("from"); from > 0 {
			q.Set("from", fmt.Sprintf("%d", from))
		}
		if to, _ := cmd.Flags().GetInt64("to"); to > 0 {
			q.Set("to", fmt.Sprintf("%d", to))
		}
		if cmd.Flags().Changed("active") {
			active, _ := cmd.Flags().GetInt("active")
			q.Set("active", fmt.Sprintf("%d", active))
		}
		data, err := client.Get(itemsPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	itemsCmd.AddCommand(itemListCmd)

	itemListCmd.Flags().Int("count", 10, "Number of items to fetch (max 100)")
	itemListCmd.Flags().Int("skip", 0, "Number of items to skip")
	itemListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch items created after this time")
	itemListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch items created before this time")
	itemListCmd.Flags().Int("active", 1, "Filter by status: 1 for active, 0 for inactive")
}
