package customers

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		name, _ := cmd.Flags().GetString("name")
		contact, _ := cmd.Flags().GetString("contact")
		email, _ := cmd.Flags().GetString("email")
		failExisting, _ := cmd.Flags().GetString("fail-existing")
		gstin, _ := cmd.Flags().GetString("gstin")
		notes, _ := cmd.Flags().GetStringArray("note")

		body := map[string]interface{}{
			"name": name,
		}
		if contact != "" {
			body["contact"] = contact
		}
		if email != "" {
			body["email"] = email
		}
		if failExisting != "" {
			body["fail_existing"] = failExisting
		}
		if gstin != "" {
			body["gstin"] = gstin
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
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

	createCmd.Flags().String("name", "", "Customer name, 3-50 characters (required)")
	createCmd.Flags().String("contact", "", "Phone number with country code (e.g. +919876543210)")
	createCmd.Flags().String("email", "", "Customer email address (max 64 characters)")
	createCmd.Flags().String("fail-existing", "", "Duplicate behaviour: 1=error (default), 0=return existing")
	createCmd.Flags().String("gstin", "", "Customer GST identification number")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	_ = createCmd.MarkFlagRequired("name")
}
