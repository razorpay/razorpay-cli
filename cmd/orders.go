package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage orders",
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all orders",
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
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			q.Set("status", status)
		}
		data, err := client.Get("/orders", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var ordersFetchCmd = &cobra.Command{
	Use:   "fetch <order_id>",
	Short: "Fetch an order by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/orders/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var ordersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new order",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		receipt, _ := cmd.Flags().GetString("receipt")
		params, _ := cmd.Flags().GetStringArray("param")

		if amount <= 0 {
			return fmt.Errorf("--amount is required and must be > 0")
		}
		if currency == "" {
			return fmt.Errorf("--currency is required")
		}

		body := map[string]interface{}{
			"amount":   amount,
			"currency": currency,
		}
		if receipt != "" {
			body["receipt"] = receipt
		}
		extra, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		for k, v := range extra {
			body[k] = v
		}

		data, err := client.Post("/orders", body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var ordersUpdateCmd = &cobra.Command{
	Use:   "update <order_id>",
	Short: "Update an order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		params, _ := cmd.Flags().GetStringArray("param")
		body, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		data, err := client.Patch("/orders/"+args[0], body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var ordersFetchPaymentsCmd = &cobra.Command{
	Use:   "payments <order_id>",
	Short: "Fetch payments for an order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/orders/"+args[0]+"/payments", nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	ordersCmd.AddCommand(ordersListCmd)
	ordersCmd.AddCommand(ordersFetchCmd)
	ordersCmd.AddCommand(ordersCreateCmd)
	ordersCmd.AddCommand(ordersUpdateCmd)
	ordersCmd.AddCommand(ordersFetchPaymentsCmd)

	ordersListCmd.Flags().Int("count", 10, "Number of orders to fetch (max 100)")
	ordersListCmd.Flags().Int("skip", 0, "Number of orders to skip")
	ordersListCmd.Flags().Int64("from", 0, "Unix timestamp: fetch orders created after this time")
	ordersListCmd.Flags().Int64("to", 0, "Unix timestamp: fetch orders created before this time")
	ordersListCmd.Flags().String("status", "", "Filter by status: created, attempted, paid")

	ordersCreateCmd.Flags().Int("amount", 0, "Order amount in smallest currency unit (e.g. paise for INR)")
	ordersCreateCmd.Flags().String("currency", "INR", "Currency code (e.g. INR)")
	ordersCreateCmd.Flags().String("receipt", "", "Receipt number for your internal reference")
	ordersCreateCmd.Flags().StringArray("param", nil, "Additional parameter as key=value")
	_ = ordersCreateCmd.MarkFlagRequired("amount")

	ordersUpdateCmd.Flags().StringArray("param", nil, "Parameter as key=value")
}
