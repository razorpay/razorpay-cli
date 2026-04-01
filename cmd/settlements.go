package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var settlementsCmd = &cobra.Command{
	Use:   "settlements",
	Short: "Manage settlements",
}

var settlementsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all settlements",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
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
		data, err := client.Get("/settlements", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var settlementsFetchCmd = &cobra.Command{
	Use:   "fetch <settlement_id>",
	Short: "Fetch a settlement by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/settlements/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var settlementsReconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Fetch settlement recon report",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		q := url.Values{}
		year, _ := cmd.Flags().GetInt("year")
		month, _ := cmd.Flags().GetInt("month")
		day, _ := cmd.Flags().GetInt("day")
		if year > 0 {
			q.Set("year", fmt.Sprintf("%d", year))
		}
		if month > 0 {
			q.Set("month", fmt.Sprintf("%d", month))
		}
		if day > 0 {
			q.Set("day", fmt.Sprintf("%d", day))
		}
		data, err := client.Get("/settlements/recon/combined", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	settlementsCmd.AddCommand(settlementsListCmd)
	settlementsCmd.AddCommand(settlementsFetchCmd)
	settlementsCmd.AddCommand(settlementsReconCmd)

	settlementsListCmd.Flags().Int("count", 10, "Number of settlements to fetch")
	settlementsListCmd.Flags().Int("skip", 0, "Number of settlements to skip")
	settlementsListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch settlements created after this time")
	settlementsListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch settlements created before this time")

	settlementsReconCmd.Flags().Int("year", 0, "Year (e.g. 2024)")
	settlementsReconCmd.Flags().Int("month", 0, "Month (1-12)")
	settlementsReconCmd.Flags().Int("day", 0, "Day (1-31)")
}
