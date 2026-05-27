package subscriptions

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/razorpay/razorpay-cli/api"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <subscription_id>",
	Short: "Update a subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		planID, _ := cmd.Flags().GetString("plan-id")
		offerID, _ := cmd.Flags().GetString("offer-id")
		quantity, _ := cmd.Flags().GetInt("quantity")
		remainingCount, _ := cmd.Flags().GetInt("remaining-count")
		startAt, _ := cmd.Flags().GetInt64("start-at")
		scheduleChangeAt, _ := cmd.Flags().GetString("schedule-change-at")
		customerNotify, _ := cmd.Flags().GetBool("customer-notify")

		body := map[string]any{}
		if planID != "" {
			body["plan_id"] = planID
		}
		if offerID != "" {
			body["offer_id"] = offerID
		}
		if quantity > 0 {
			body["quantity"] = quantity
		}
		if remainingCount > 0 {
			body["remaining_count"] = remainingCount
		}
		if startAt > 0 {
			body["start_at"] = startAt
		}
		if scheduleChangeAt != "" {
			body["schedule_change_at"] = scheduleChangeAt
		}
		if cmd.Flags().Changed("customer-notify") {
			body["customer_notify"] = customerNotify
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
		}

		data, err := client.Patch(basePath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)

	updateCmd.Flags().String("plan-id", "", "New plan ID to link to the subscription")
	updateCmd.Flags().String("offer-id", "", "Offer ID to link to the subscription")
	updateCmd.Flags().Int("quantity", 0, "Number of times the plan amount is charged per invoice")
	updateCmd.Flags().Int("remaining-count", 0, "Override the remaining billing cycles count")
	updateCmd.Flags().Int64("start-at", 0, "New start date as Unix timestamp")
	updateCmd.Flags().String("schedule-change-at", "", "When to apply the update: now (default) or cycle_end")
	updateCmd.Flags().Bool("customer-notify", true, "Let Razorpay send update notifications to customer")
}
