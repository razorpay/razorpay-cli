package invoices

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var itemUpdateCmd = &cobra.Command{
	Use:   "update <item_id>",
	Short: "Update an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		active, _ := cmd.Flags().GetBool("active")

		body := map[string]any{}
		if name != "" {
			body["name"] = name
		}
		if description != "" {
			body["description"] = description
		}
		if amount > 0 {
			body["amount"] = amount
		}
		if currency != "" {
			body["currency"] = currency
		}
		if cmd.Flags().Changed("active") {
			body["active"] = active
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
		}

		data, err := client.Patch(itemsPath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	itemsCmd.AddCommand(itemUpdateCmd)

	itemUpdateCmd.Flags().String("name", "", "Item name")
	itemUpdateCmd.Flags().String("description", "", "Brief description of the item")
	itemUpdateCmd.Flags().Int64("amount", 0, "Item price in smallest currency unit e.g. paise")
	itemUpdateCmd.Flags().String("currency", "", "Currency code e.g. INR")
	itemUpdateCmd.Flags().Bool("active", true, "Item status: true for active, false to deactivate")
}
