package payments

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <payment_id>",
	Short: "Update notes on a payment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		notes, _ := cmd.Flags().GetStringArray("note")

		if len(notes) == 0 {
			return fmt.Errorf("--note is required (only notes can be updated on a payment)")
		}

		notesMap, err := api.ParseParams(notes)
		if err != nil {
			return err
		}

		body := map[string]any{
			"notes": notesMap,
		}

		data, err := client.Patch(basePath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(updateCmd)

	updateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
