package smartcollect

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var tpvCreateCmd = &cobra.Command{
	Use:   "tpv-create",
	Short: "Create a TPV (Third-Party Validation) virtual account with allowed payers",
	Long: `Create a TPV virtual account with allowed payers.

Pass --allowed-payers as a JSON array. Each element supports: type, bank_account.ifsc, bank_account.account_number

Example:
  razorpay smart-collect tpv-create \
    --receiver-type bank_account \
    --bank-account-descriptor "ACME Corp" \
    --allowed-payers '[
      {"type":"bank_account","bank_account":{"ifsc":"UTIB0000001","account_number":"9876543210123456"}},
      {"type":"bank_account","bank_account":{"ifsc":"HDFC0000002","account_number":"1234567890123456"}}
    ]'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		receiverTypes, _ := cmd.Flags().GetStringArray("receiver-type")
		bankDescriptor, _ := cmd.Flags().GetString("bank-account-descriptor")
		allowedPayersJSON, _ := cmd.Flags().GetString("allowed-payers")

		if len(receiverTypes) == 0 {
			return fmt.Errorf("at least one --receiver-type is required (e.g. --receiver-type bank_account)")
		}
		if allowedPayersJSON == "" {
			return fmt.Errorf("--allowed-payers is required (JSON array of payer objects)")
		}

		var allowedPayers any
		if err := json.Unmarshal([]byte(allowedPayersJSON), &allowedPayers); err != nil {
			return fmt.Errorf("--allowed-payers is not valid JSON: %w", err)
		}

		receivers := map[string]any{
			"types": receiverTypes,
		}
		if bankDescriptor != "" {
			receivers["bank_account"] = map[string]any{
				"descriptor": bankDescriptor,
			}
		}

		body := map[string]any{
			"receivers":      receivers,
			"allowed_payers": allowedPayers,
		}

		description, _ := cmd.Flags().GetString("description")
		customerID, _ := cmd.Flags().GetString("customer-id")
		closeBy, _ := cmd.Flags().GetInt64("close-by")
		notes, _ := cmd.Flags().GetStringArray("note")

		if description != "" {
			body["description"] = description
		}
		if customerID != "" {
			body["customer_id"] = customerID
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

	tpvCreateCmd.Flags().StringArray("receiver-type", nil, "Receiver type (repeatable, e.g. --receiver-type bank_account)")
	tpvCreateCmd.Flags().String("bank-account-descriptor", "", "Bank account descriptor (identifier label for the bank account receiver)")
	tpvCreateCmd.Flags().String("allowed-payers", "", `Allowed payers as a JSON array (required, max 10). Each object: {"type":"bank_account","bank_account":{"ifsc":"...","account_number":"..."}}`)
	tpvCreateCmd.Flags().String("description", "", "Description of the virtual account")
	tpvCreateCmd.Flags().String("customer-id", "", "Customer ID to associate with this virtual account")
	tpvCreateCmd.Flags().Int64("close-by", 0, "Unix timestamp for automatic closure (min 15 mins from now)")
	tpvCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable)")
}
