package customers

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <customer_id>",
	Short: "Update a customer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		contact, _ := cmd.Flags().GetString("contact")

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if email != "" {
			body["email"] = email
		}
		if contact != "" {
			body["contact"] = contact
		}

		data, err := client.Put(basePath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)

	updateCmd.Flags().String("name", "", "Customer name, 3-50 characters")
	updateCmd.Flags().String("email", "", "Customer email address (max 64 characters)")
	updateCmd.Flags().String("contact", "", "Phone number with country code (e.g. +919876543210)")
}
