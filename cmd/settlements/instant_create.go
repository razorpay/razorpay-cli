package settlements

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var instantCreateCmd = &cobra.Command{
	Use:   "instant-create",
	Short: "Create an instant (on-demand) settlement",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		amount, _ := cmd.Flags().GetInt64("amount")
		settleFullBalance, _ := cmd.Flags().GetBool("settle-full-balance")
		description, _ := cmd.Flags().GetString("description")
		notes, _ := cmd.Flags().GetStringArray("note")

		if !settleFullBalance && amount <= 0 {
			return fmt.Errorf("--amount is required unless --settle-full-balance is set")
		}

		body := map[string]any{
			"settle_full_balance": settleFullBalance,
		}
		if amount > 0 {
			body["amount"] = amount
		}
		if description != "" {
			body["description"] = description
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}
		data, err := client.Post(ondemandPath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(instantCreateCmd)

	instantCreateCmd.Flags().Int64("amount", 0, "Amount to settle in paise (required unless --settle-full-balance)")
	instantCreateCmd.Flags().Bool("settle-full-balance", false, "Settle the maximum possible amount ignoring --amount")
	instantCreateCmd.Flags().String("description", "", "Custom note for the settlement (max 30 chars)")
	instantCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
