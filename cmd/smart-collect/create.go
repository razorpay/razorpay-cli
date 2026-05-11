package smartcollect

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a virtual account (customer identifier)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		receiverTypes, _ := cmd.Flags().GetStringArray("receiver-type")
		bankDescriptor, _ := cmd.Flags().GetString("bank-account-descriptor")
		vpaDescriptor, _ := cmd.Flags().GetString("vpa-descriptor")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		customerID, _ := cmd.Flags().GetString("customer-id")
		amountExpected, _ := cmd.Flags().GetInt64("amount-expected")
		closeBy, _ := cmd.Flags().GetInt64("close-by")
		notes, _ := cmd.Flags().GetStringArray("note")

		if len(receiverTypes) == 0 {
			return fmt.Errorf("at least one --receiver-type is required (bank_account or vpa)")
		}

		receivers := map[string]any{
			"types": receiverTypes,
		}
		if bankDescriptor != "" {
			receivers["bank_account"] = map[string]any{
				"descriptor": bankDescriptor,
			}
		}
		if vpaDescriptor != "" {
			receivers["vpa"] = map[string]any{
				"descriptor": vpaDescriptor,
			}
		}
		body := map[string]any{
			"receivers": receivers,
		}
		if name != "" {
			body["name"] = name
		}
		if description != "" {
			body["description"] = description
		}
		if customerID != "" {
			body["customer_id"] = customerID
		}
		if amountExpected > 0 {
			body["amount_expected"] = amountExpected
		}
		if closeBy > 0 {
			body["close_by"] = closeBy
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

	createCmd.Flags().StringArray("receiver-type", nil, "Receiver type: bank_account or vpa (required, repeatable)")
	createCmd.Flags().String("bank-account-descriptor", "", "Bank account descriptor (identifier label for the bank account receiver)")
	createCmd.Flags().String("vpa-descriptor", "", "VPA descriptor (identifier label for the VPA receiver)")
	createCmd.Flags().String("name", "", "Name for the virtual account")
	createCmd.Flags().String("description", "", "Description of the virtual account")
	createCmd.Flags().String("customer-id", "", "Customer ID to associate with this virtual account")
	createCmd.Flags().Int64("amount-expected", 0, "Expected payment amount in paise (0 = any amount)")
	createCmd.Flags().Int64("close-by", 0, "Unix timestamp for automatic closure (min 15 minutes from now)")
	createCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
