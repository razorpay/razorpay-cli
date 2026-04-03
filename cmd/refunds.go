package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var refundsCmd = &cobra.Command{
	Use:   "refunds",
	Short: "Manage refunds",
}

var refundsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all refunds",
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
		data, err := client.Get("/refunds", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var refundsFetchCmd = &cobra.Command{
	Use:   "fetch <refund_id>",
	Short: "Fetch a refund by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/refunds/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var refundsCreateCmd = &cobra.Command{
	Use:   "create <payment_id>",
	Short: "Create a refund for a payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		amount, _ := cmd.Flags().GetInt("amount")
		speed, _ := cmd.Flags().GetString("speed")
		params, _ := cmd.Flags().GetStringArray("param")

		body := map[string]interface{}{}
		if amount > 0 {
			body["amount"] = amount
		}
		if speed != "" {
			body["speed"] = speed
		}
		extra, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		for k, v := range extra {
			body[k] = v
		}

		data, err := client.Post("/payments/"+args[0]+"/refund", body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var refundsUpdateCmd = &cobra.Command{
	Use:   "update <refund_id>",
	Short: "Update a refund",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		params, _ := cmd.Flags().GetStringArray("param")
		body, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		data, err := client.Patch("/refunds/"+args[0], body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	refundsCmd.AddCommand(refundsListCmd)
	refundsCmd.AddCommand(refundsFetchCmd)
	refundsCmd.AddCommand(refundsCreateCmd)
	refundsCmd.AddCommand(refundsUpdateCmd)

	refundsListCmd.Flags().Int("count", 10, "Number of refunds to fetch")
	refundsListCmd.Flags().Int("skip", 0, "Number of refunds to skip")
	refundsListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch refunds created after this time")
	refundsListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch refunds created before this time")

	refundsCreateCmd.Flags().Int("amount", 0, "Amount to refund in smallest currency unit (omit for full refund)")
	refundsCreateCmd.Flags().String("speed", "", "Refund speed: normal or optimum")
	refundsCreateCmd.Flags().StringArray("param", nil, "Additional parameter as key=value")

	refundsUpdateCmd.Flags().StringArray("param", nil, "Parameter as key=value")
}
