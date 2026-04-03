package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var disputesCmd = &cobra.Command{
	Use:   "disputes",
	Short: "Manage disputes",
}

var disputesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all disputes",
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
		data, err := client.Get("/disputes", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var disputesFetchCmd = &cobra.Command{
	Use:   "fetch <dispute_id>",
	Short: "Fetch a dispute by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/disputes/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var disputesAcceptCmd = &cobra.Command{
	Use:   "accept <dispute_id>",
	Short: "Accept a dispute",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Post("/disputes/"+args[0]+"/accept", nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var disputesContestCmd = &cobra.Command{
	Use:   "contest <dispute_id>",
	Short: "Contest a dispute",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		params, _ := cmd.Flags().GetStringArray("param")
		body, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		action, _ := cmd.Flags().GetString("action")
		if action != "" {
			body["action"] = action
		}

		var path string
		if action == "submit" {
			path = fmt.Sprintf("/disputes/%s/contest", args[0])
		} else {
			path = fmt.Sprintf("/disputes/%s/contest", args[0])
		}
		data, err := client.Patch(path, body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	disputesCmd.AddCommand(disputesListCmd)
	disputesCmd.AddCommand(disputesFetchCmd)
	disputesCmd.AddCommand(disputesAcceptCmd)
	disputesCmd.AddCommand(disputesContestCmd)

	disputesListCmd.Flags().Int("count", 10, "Number of disputes to fetch")
	disputesListCmd.Flags().Int("skip", 0, "Number of disputes to skip")
	disputesListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch disputes created after this time")
	disputesListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch disputes created before this time")

	disputesContestCmd.Flags().StringArray("param", nil, "Parameter as key=value (e.g. --param amount=1000)")
	disputesContestCmd.Flags().String("action", "", "Action: draft or submit")
}
