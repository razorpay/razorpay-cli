package invoices

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var itemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new item",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		name, _ := cmd.Flags().GetString("name")
		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if amount <= 0 {
			return fmt.Errorf("--amount is required and must be > 0")
		}
		if currency == "" {
			return fmt.Errorf("--currency is required")
		}

		body := map[string]any{
			"name":     name,
			"amount":   amount,
			"currency": currency,
		}
		if description != "" {
			body["description"] = description
		}

		data, err := client.Post(itemsPath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	itemsCmd.AddCommand(itemCreateCmd)

	itemCreateCmd.Flags().String("name", "", "Item name (required)")
	itemCreateCmd.Flags().Int64("amount", 0, "Item price in smallest currency unit e.g. paise (required)")
	itemCreateCmd.Flags().String("currency", "", "Currency code e.g. INR (required)")
	itemCreateCmd.Flags().String("description", "", "Brief description of the item")
}
