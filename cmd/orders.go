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
	Long:  "Create, list, fetch, and update Razorpay orders.",
}

var ordersListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List orders",
	Long:    "List orders on the account, with optional pagination, status, and a created-at time window.",
	Example: "  razorpay orders list --status paid --count 25",
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
	Use:     "fetch <order_id>",
	Short:   "Fetch an order by ID",
	Example: "  razorpay orders fetch order_DBJOWzybf0sJbb",
	Args:    cobra.ExactArgs(1),
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
	Use:     "create",
	Short:   "Create an order",
	Example: "  razorpay orders create --amount 50000 --currency INR --receipt rcpt_001",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		receipt, _ := cmd.Flags().GetString("receipt")
		params, _ := cmd.Flags().GetStringArray("param")

		if amount <= 0 {
			return fmt.Errorf("amount must be greater than 0")
		}
		if currency == "" {
			return fmt.Errorf("currency is required")
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
	Use:     "update <order_id>",
	Short:   "Update an order's notes",
	Example: "  razorpay orders update order_DBJOWzybf0sJbb --param notes[shipment]=AWB1234",
	Args:    cobra.ExactArgs(1),
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
	Short: "List payments for an order",
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

	ordersListCmd.Flags().Int("count", 10, "Maximum number of orders to return (max 100)")
	ordersListCmd.Flags().Int("skip", 0, "Number of orders to skip for pagination")
	ordersListCmd.Flags().Int64("from", 0, "Include orders created on or after this Unix timestamp")
	ordersListCmd.Flags().Int64("to", 0, "Include orders created on or before this Unix timestamp")
	ordersListCmd.Flags().String("status", "", "Filter by order status: created, attempted, or paid")

	ordersCreateCmd.Flags().Int("amount", 0, "Order amount in the smallest currency unit (e.g. paise for INR)")
	ordersCreateCmd.Flags().String("currency", "INR", "ISO 4217 currency code (e.g. INR)")
	ordersCreateCmd.Flags().String("receipt", "", "Receipt number for your internal reference")
	ordersCreateCmd.Flags().StringArray("param", nil, "Additional field as key=value; repeatable")
	_ = ordersCreateCmd.MarkFlagRequired("amount")

	ordersUpdateCmd.Flags().StringArray("param", nil, "Field to update as key=value; repeatable")
}
