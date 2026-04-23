package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a direct transfer to a linked account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		account, _ := cmd.Flags().GetString("account")
		amount, _ := cmd.Flags().GetInt64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		onHoldUntil, _ := cmd.Flags().GetInt64("on-hold-until")
		notes, _ := cmd.Flags().GetStringArray("note")
		linkedAccountNotes, _ := cmd.Flags().GetStringArray("linked-account-note")

		if account == "" {
			return fmt.Errorf("--account is required (linked account ID)")
		}
		if amount == 0 {
			return fmt.Errorf("--amount is required")
		}

		body := map[string]any{
			"account":  account,
			"amount":   amount,
			"currency": currency,
		}
		if cmd.Flags().Changed("on-hold") {
			onHold, _ := cmd.Flags().GetBool("on-hold")
			body["on_hold"] = onHold
		}
		if onHoldUntil > 0 {
			body["on_hold_until"] = onHoldUntil
		}
		if len(linkedAccountNotes) > 0 {
			body["linked_account_notes"] = linkedAccountNotes
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		var data []byte
		var err error
		if key, _ := cmd.Flags().GetString("idempotency-key"); key != "" {
			data, err = client.PostWithHeaders(transfersPath, body, map[string]string{
				"X-Transfer-Idempotency": key,
			})
		} else {
			data, err = client.Post(transfersPath, body)
		}
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferCreateCmd)

	transferCreateCmd.Flags().String("account", "", "Linked account ID to transfer to (required)")
	transferCreateCmd.Flags().Int64("amount", 0, "Transfer amount in paise (required, min 100)")
	transferCreateCmd.Flags().String("currency", "INR", "Currency code (default INR)")
	transferCreateCmd.Flags().Bool("on-hold", false, "Put the transfer settlement on hold")
	transferCreateCmd.Flags().Int64("on-hold-until", 0, "Unix timestamp until which settlement is held")
	transferCreateCmd.Flags().StringArray("linked-account-note", nil, "Note key to expose to the linked account (repeatable)")
	transferCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
	transferCreateCmd.Flags().String("idempotency-key", "", "Idempotency key (4-36 chars: alphanumerics, hyphens, underscores, spaces) to safely retry the transfer")
}
