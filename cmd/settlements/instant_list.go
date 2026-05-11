package settlements

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var instantListCmd = &cobra.Command{
	Use:   "instant-list",
	Short: "List all instant (on-demand) settlements",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		if expands, _ := cmd.Flags().GetStringArray("expand"); len(expands) > 0 {
			for _, e := range expands {
				q.Add("expand[]", e)
			}
		}
		data, err := client.Get(ondemandPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(instantListCmd)

	instantListCmd.Flags().Int("count", 10, "Number of instant settlements to fetch (max 100)")
	instantListCmd.Flags().Int("skip", 0, "Number of instant settlements to skip")
	instantListCmd.Flags().StringArray("expand", nil, "Expand related objects (e.g. --expand ondemand_payouts)")
}
