package settlements

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Fetch settlement recon report",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}

		for _, flag := range []string{"year", "month", "day"} {
			if s, _ := cmd.Flags().GetString(flag); s != "" {
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil || n <= 0 {
					return fmt.Errorf("--%s: %q is not a valid number", flag, s)
				}
				q.Set(flag, fmt.Sprintf("%d", n))
			}
		}
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}

		data, err := client.Get(reconPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(reconCmd)

	reconCmd.Flags().String("year", "", "Year (e.g. 2025)")
	reconCmd.Flags().String("month", "", "Month (1-12)")
	reconCmd.Flags().String("day", "", "Day (1-31)")
	reconCmd.Flags().Int("count", 10, "Number of recon records to fetch (max 1000)")
	reconCmd.Flags().Int("skip", 0, "Number of recon records to skip")
}
