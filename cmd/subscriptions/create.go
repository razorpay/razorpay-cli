package subscriptions

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new subscription (or subscription link)",
	Long: `Create a new subscription.

To send the subscription link to a customer via Razorpay, pass
--notify-info-email and/or --notify-info-phone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		planID, _ := cmd.Flags().GetString("plan-id")
		totalCount, _ := cmd.Flags().GetInt("total-count")
		quantity, _ := cmd.Flags().GetInt("quantity")
		startAt, _ := cmd.Flags().GetInt64("start-at")
		expireBy, _ := cmd.Flags().GetInt64("expire-by")
		customerNotify, _ := cmd.Flags().GetBool("customer-notify")
		addonsJSON, _ := cmd.Flags().GetString("addons")
		offerID, _ := cmd.Flags().GetString("offer-id")
		notes, _ := cmd.Flags().GetStringArray("note")
		notifyEmail, _ := cmd.Flags().GetString("notify-info-email")
		notifyPhone, _ := cmd.Flags().GetString("notify-info-phone")

		if planID == "" {
			return fmt.Errorf("--plan-id is required")
		}
		if totalCount == 0 {
			return fmt.Errorf("--total-count is required")
		}

		body := map[string]any{
			"plan_id":     planID,
			"total_count": totalCount,
		}
		if quantity > 0 {
			body["quantity"] = quantity
		}
		if startAt > 0 {
			body["start_at"] = startAt
		}
		if expireBy > 0 {
			body["expire_by"] = expireBy
		}
		if cmd.Flags().Changed("customer-notify") {
			body["customer_notify"] = customerNotify
		}
		if offerID != "" {
			body["offer_id"] = offerID
		}
		if addonsJSON != "" {
			var addons any
			if err := json.Unmarshal([]byte(addonsJSON), &addons); err != nil {
				return fmt.Errorf("--addons is not valid JSON: %w", err)
			}
			body["addons"] = addons
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}
		notifyInfo := map[string]any{}
		if notifyEmail != "" {
			notifyInfo["notify_email"] = notifyEmail
		}
		if notifyPhone != "" {
			notifyInfo["notify_phone"] = notifyPhone
		}
		if len(notifyInfo) > 0 {
			body["notify_info"] = notifyInfo
		}

		data, err := client.Post(basePath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(createCmd)

	createCmd.Flags().String("plan-id", "", "Plan ID to subscribe to (required)")
	createCmd.Flags().Int("total-count", 0, "Number of billing cycles (required)")
	createCmd.Flags().Int("quantity", 1, "Times to charge the plan amount per invoice (default 1)")
	createCmd.Flags().Int64("start-at", 0, "Subscription start as Unix timestamp (default: immediately after auth)")
	createCmd.Flags().Int64("expire-by", 0, "Authorization payment deadline as Unix timestamp")
	createCmd.Flags().Bool("customer-notify", true, "Let Razorpay send notifications to customer (default true)")
	createCmd.Flags().String("offer-id", "", "Offer ID to link to the subscription")
	createCmd.Flags().String("addons", "", "Upfront addon charges as a JSON array")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	createCmd.Flags().String("notify-info-email", "", "Customer email for subscription link delivery")
	createCmd.Flags().String("notify-info-phone", "", "Customer phone for subscription link delivery")
}
