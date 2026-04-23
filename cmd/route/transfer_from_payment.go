package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFromPaymentCmd = &cobra.Command{
	Use:   "create-from-payment <payment_id>",
	Short: "Create transfers from a captured payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		account, _ := cmd.Flags().GetString("account")
		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		onHoldUntil, _ := cmd.Flags().GetInt64("on-hold-until")
		linkedAccountNotes, _ := cmd.Flags().GetStringArray("linked-account-note")
		notes, _ := cmd.Flags().GetStringArray("note")

		if account == "" {
			return fmt.Errorf("--account is required")
		}
		if amount == 0 {
			return fmt.Errorf("--amount is required")
		}

		transfer := map[string]any{
			"account":  account,
			"amount":   amount,
			"currency": currency,
		}
		if cmd.Flags().Changed("on-hold") {
			onHold, _ := cmd.Flags().GetBool("on-hold")
			transfer["on_hold"] = onHold
		}
		if onHoldUntil > 0 {
			transfer["on_hold_until"] = onHoldUntil
		}
		if len(linkedAccountNotes) > 0 {
			transfer["linked_account_notes"] = linkedAccountNotes
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			transfer["notes"] = notesMap
		}

		body := map[string]any{
			"transfers": []any{transfer},
		}

		data, err := client.Post("/payments/"+args[0]+"/transfers", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFromPaymentCmd)

	transferFromPaymentCmd.Flags().String("account", "", "Linked account ID to transfer to (required)")
	transferFromPaymentCmd.Flags().Int64("amount", 0, "Transfer amount in paise (required)")
	transferFromPaymentCmd.Flags().String("currency", "INR", "Currency code (default INR)")
	transferFromPaymentCmd.Flags().Bool("on-hold", false, "Put the transfer settlement on hold")
	transferFromPaymentCmd.Flags().Int64("on-hold-until", 0, "Unix timestamp until which settlement is held")
	transferFromPaymentCmd.Flags().StringArray("linked-account-note", nil, "Note key to expose to the linked account (repeatable)")
	transferFromPaymentCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
