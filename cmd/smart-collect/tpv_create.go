package smartcollect

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var tpvCreateCmd = &cobra.Command{
	Use:   "tpv-create",
	Short: "Create a TPV (Third-Party Validation) virtual account with allowed payers",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		ifscList, _ := cmd.Flags().GetStringArray("ifsc")
		accountNumberList, _ := cmd.Flags().GetStringArray("account-number")

		if len(ifscList) == 0 || len(accountNumberList) == 0 {
			return fmt.Errorf("at least one --ifsc and --account-number pair is required")
		}
		if len(ifscList) != len(accountNumberList) {
			return fmt.Errorf("--ifsc and --account-number must be provided in equal numbers (one pair per payer)")
		}
		if len(ifscList) > 10 {
			return fmt.Errorf("maximum 10 allowed payers permitted")
		}

		allowedPayers := make([]map[string]any, len(ifscList))
		for i := range ifscList {
			allowedPayers[i] = map[string]any{
				"type": "bank_account",
				"bank_account": map[string]any{
					"ifsc":           ifscList[i],
					"account_number": accountNumberList[i],
				},
			}
		}

		body := map[string]any{
			"receivers": map[string]any{
				"types": []string{"bank_account"},
			},
			"allowed_payers": allowedPayers,
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		customerID, _ := cmd.Flags().GetString("customer-id")
		amountExpected, _ := cmd.Flags().GetInt64("amount-expected")
		closeBy, _ := cmd.Flags().GetInt64("close-by")
		notes, _ := cmd.Flags().GetStringArray("note")

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
	Cmd.AddCommand(tpvCreateCmd)

	tpvCreateCmd.Flags().StringArray("ifsc", nil, "Bank IFSC code for an allowed payer (repeatable, one per payer)")
	tpvCreateCmd.Flags().StringArray("account-number", nil, "Bank account number for an allowed payer (repeatable, one per payer)")
	tpvCreateCmd.Flags().String("name", "", "Name for the virtual account")
	tpvCreateCmd.Flags().String("description", "", "Description of the virtual account")
	tpvCreateCmd.Flags().String("customer-id", "", "Customer ID to associate with this virtual account")
	tpvCreateCmd.Flags().Int64("amount-expected", 0, "Expected payment amount in paise")
	tpvCreateCmd.Flags().Int64("close-by", 0, "Unix timestamp for automatic closure")
	tpvCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable)")
}
