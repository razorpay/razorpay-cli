package subscriptions

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var planCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		period, _ := cmd.Flags().GetString("period")
		interval, _ := cmd.Flags().GetInt("interval")
		itemName, _ := cmd.Flags().GetString("item-name")
		itemAmount, _ := cmd.Flags().GetInt64("item-amount")
		itemCurrency, _ := cmd.Flags().GetString("item-currency")
		itemDescription, _ := cmd.Flags().GetString("item-description")
		notes, _ := cmd.Flags().GetStringArray("note")

		if period == "" {
			return fmt.Errorf("--period is required (daily, weekly, monthly, quarterly, yearly)")
		}
		if interval == 0 {
			return fmt.Errorf("--interval is required")
		}
		if itemName == "" {
			return fmt.Errorf("--item-name is required")
		}
		if itemAmount <= 0 {
			return fmt.Errorf("--item-amount is required and must be > 0")
		}

		item := map[string]any{
			"name":     itemName,
			"amount":   itemAmount,
			"currency": itemCurrency,
		}
		if itemDescription != "" {
			item["description"] = itemDescription
		}

		body := map[string]any{
			"period":   period,
			"interval": interval,
			"item":     item,
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		data, err := client.Post(plansPath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	plansCmd.AddCommand(planCreateCmd)

	planCreateCmd.Flags().String("period", "", "Billing period: daily, weekly, monthly, quarterly, or yearly (required)")
	planCreateCmd.Flags().Int("interval", 0, "Billing frequency multiplier, e.g. 3 for every 3 months (required)")
	planCreateCmd.Flags().String("item-name", "", "Plan name (required)")
	planCreateCmd.Flags().Int64("item-amount", 0, "Amount per billing cycle in paise (required)")
	planCreateCmd.Flags().String("item-currency", "INR", "Currency code (default INR)")
	planCreateCmd.Flags().String("item-description", "", "Plan description")
	planCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
