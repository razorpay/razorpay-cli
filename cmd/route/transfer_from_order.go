package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFromOrderCmd = &cobra.Command{
	Use:   "create-from-order",
	Short: "Create an order with embedded transfers to linked accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		receipt, _ := cmd.Flags().GetString("receipt")
		account, _ := cmd.Flags().GetString("account")
		transferAmount, _ := cmd.Flags().GetInt64("transfer-amount")
		transferCurrency, _ := cmd.Flags().GetString("transfer-currency")
		onHoldUntil, _ := cmd.Flags().GetInt64("on-hold-until")
		linkedAccountNotes, _ := cmd.Flags().GetStringArray("linked-account-note")
		notes, _ := cmd.Flags().GetStringArray("note")

		if amount == 0 {
			return fmt.Errorf("--amount is required")
		}
		if account == "" {
			return fmt.Errorf("--account is required")
		}
		if transferAmount == 0 {
			return fmt.Errorf("--transfer-amount is required")
		}

		transfer := map[string]any{
			"account":  account,
			"amount":   transferAmount,
			"currency": transferCurrency,
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
			"amount":    amount,
			"currency":  currency,
			"transfers": []any{transfer},
		}
		if receipt != "" {
			body["receipt"] = receipt
		}

		data, err := client.Post("/orders", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFromOrderCmd)

	transferFromOrderCmd.Flags().Int64("amount", 0, "Order amount in paise (required)")
	transferFromOrderCmd.Flags().String("currency", "INR", "Order currency code (default INR)")
	transferFromOrderCmd.Flags().String("receipt", "", "Order receipt / reference identifier")
	transferFromOrderCmd.Flags().String("account", "", "Linked account ID to transfer to (required)")
	transferFromOrderCmd.Flags().Int64("transfer-amount", 0, "Transfer amount in paise (required)")
	transferFromOrderCmd.Flags().String("transfer-currency", "INR", "Transfer currency code (default INR)")
	transferFromOrderCmd.Flags().Bool("on-hold", false, "Put the transfer settlement on hold")
	transferFromOrderCmd.Flags().Int64("on-hold-until", 0, "Unix timestamp until which settlement is held")
	transferFromOrderCmd.Flags().StringArray("linked-account-note", nil, "Note key to expose to the linked account (repeatable)")
	transferFromOrderCmd.Flags().StringArray("note", nil, "Note as key=value for the transfer (repeatable, max 15 pairs)")
}
