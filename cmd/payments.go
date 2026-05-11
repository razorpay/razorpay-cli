package cmd

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage payments",
	Long:  "List, fetch, capture, and update Razorpay payments.",
}

var paymentsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List payments",
	Long:    "List payments on the account, with optional pagination and a created-at time window.",
	Example: "  razorpay payments list --count 25 --from 1704067200",
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
		data, err := client.Get("/payments", q)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var paymentsFetchCmd = &cobra.Command{
	Use:     "fetch <payment_id>",
	Short:   "Fetch a payment by ID",
	Example: "  razorpay payments fetch pay_29QQoUBi66xm2f",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/payments/"+args[0], nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var paymentsCaptureCmd = &cobra.Command{
	Use:     "capture <payment_id>",
	Short:   "Capture an authorized payment",
	Example: "  razorpay payments capture pay_29QQoUBi66xm2f --amount 50000 --currency INR",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		amount, _ := cmd.Flags().GetInt("amount")
		currency, _ := cmd.Flags().GetString("currency")
		if amount <= 0 {
			return fmt.Errorf("amount must be greater than 0")
		}
		body := map[string]interface{}{
			"amount":   amount,
			"currency": currency,
		}
		data, err := client.Post("/payments/"+args[0]+"/capture", body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var paymentsUpdateCmd = &cobra.Command{
	Use:     "update <payment_id>",
	Short:   "Update a payment's notes",
	Example: "  razorpay payments update pay_29QQoUBi66xm2f --param notes[reason]=duplicate",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		params, _ := cmd.Flags().GetStringArray("param")
		body, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		data, err := client.Patch("/payments/"+args[0], body)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

var paymentsFetchTransfersCmd = &cobra.Command{
	Use:   "transfers <payment_id>",
	Short: "List transfers for a payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		data, err := client.Get("/payments/"+args[0]+"/transfers", nil)
		if err != nil {
			handleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	paymentsCmd.AddCommand(paymentsListCmd)
	paymentsCmd.AddCommand(paymentsFetchCmd)
	paymentsCmd.AddCommand(paymentsCaptureCmd)
	paymentsCmd.AddCommand(paymentsUpdateCmd)
	paymentsCmd.AddCommand(paymentsFetchTransfersCmd)

	paymentsListCmd.Flags().Int("count", 10, "Maximum number of payments to return (max 100)")
	paymentsListCmd.Flags().Int("skip", 0, "Number of payments to skip for pagination")
	paymentsListCmd.Flags().Int64("from", 0, "Include payments created on or after this Unix timestamp")
	paymentsListCmd.Flags().Int64("to", 0, "Include payments created on or before this Unix timestamp")

	paymentsCaptureCmd.Flags().Int("amount", 0, "Amount to capture, in the smallest currency unit (e.g. paise for INR)")
	paymentsCaptureCmd.Flags().String("currency", "INR", "ISO 4217 currency code (e.g. INR)")
	_ = paymentsCaptureCmd.MarkFlagRequired("amount")

	paymentsUpdateCmd.Flags().StringArray("param", nil, "Field to update as key=value; repeatable (e.g. --param notes[reason]=duplicate)")
}
