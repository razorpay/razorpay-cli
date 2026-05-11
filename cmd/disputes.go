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
	Long:  "List, fetch, accept, and contest Razorpay disputes.",
}

var disputesListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List disputes",
	Example: "  razorpay disputes list --count 25",
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
	Use:     "fetch <dispute_id>",
	Short:   "Fetch a dispute by ID",
	Example: "  razorpay disputes fetch disp_FmnsM5fHkRGQAk",
	Args:    cobra.ExactArgs(1),
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
	Use:     "accept <dispute_id>",
	Short:   "Accept a dispute",
	Long:    "Accept the dispute. The disputed amount is debited from your account and refunded to the customer.",
	Example: "  razorpay disputes accept disp_FmnsM5fHkRGQAk",
	Args:    cobra.ExactArgs(1),
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
	Use:     "contest <dispute_id>",
	Short:   "Contest a dispute",
	Long:    "Contest a dispute by attaching evidence. Use --action draft to save evidence without submitting, or --action submit to submit it to Razorpay.",
	Example: "  razorpay disputes contest disp_FmnsM5fHkRGQAk --action submit --param evidence[summary]=\"Item delivered\"",
	Args:    cobra.ExactArgs(1),
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

		path := fmt.Sprintf("/disputes/%s/contest", args[0])
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

	disputesListCmd.Flags().Int("count", 10, "Maximum number of disputes to return (max 100)")
	disputesListCmd.Flags().Int("skip", 0, "Number of disputes to skip for pagination")
	disputesListCmd.Flags().Int64("from", 0, "Include disputes created on or after this Unix timestamp")
	disputesListCmd.Flags().Int64("to", 0, "Include disputes created on or before this Unix timestamp")

	disputesContestCmd.Flags().StringArray("param", nil, "Evidence field as key=value; repeatable (e.g. --param evidence[summary]=\"Item delivered\")")
	disputesContestCmd.Flags().String("action", "", "Contest action: draft (save) or submit (finalize)")
}
