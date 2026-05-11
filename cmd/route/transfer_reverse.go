package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferReverseCmd = &cobra.Command{
	Use:   "reverse <transfer_id>",
	Short: "Reverse a transfer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		if amount, _ := cmd.Flags().GetInt64("amount"); amount > 0 {
			body["amount"] = amount
		}
		if notes, _ := cmd.Flags().GetStringArray("note"); len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		data, err := client.Post(transfersPath+"/"+args[0]+"/reversals", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferReverseCmd)

	transferReverseCmd.Flags().Int64("amount", 0, "Amount to reverse in paise (omit to reverse full amount)")
	transferReverseCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
