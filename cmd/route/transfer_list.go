package route

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transfers",
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
		if expandSettlement, _ := cmd.Flags().GetBool("expand-settlement"); expandSettlement {
			q.Set("expand[]", "recipient_settlement")
		}
		data, err := client.Get(transfersPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferListCmd)

	transferListCmd.Flags().Int("count", 10, "Number of transfers to fetch (max 100)")
	transferListCmd.Flags().Int("skip", 0, "Number of transfers to skip")
	transferListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch transfers created after this time")
	transferListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch transfers created before this time")
	transferListCmd.Flags().Bool("expand-settlement", false, "Include recipient_settlement details in response")
}
