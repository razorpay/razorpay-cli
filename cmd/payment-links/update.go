package paymentlinks

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <payment_link_id>",
	Short: "Update a payment link (only in created or partially_paid state)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		acceptPartial, _ := cmd.Flags().GetBool("accept-partial")
		referenceID, _ := cmd.Flags().GetString("reference-id")
		expireBy, _ := cmd.Flags().GetInt64("expire-by")
		notes, _ := cmd.Flags().GetStringArray("note")

		body := map[string]interface{}{}
		if cmd.Flags().Changed("accept-partial") {
			body["accept_partial"] = acceptPartial
		}
		if referenceID != "" {
			body["reference_id"] = referenceID
		}
		if expireBy > 0 {
			body["expire_by"] = expireBy
		}
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
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

	updateCmd.Flags().Bool("accept-partial", false, "Allow customer to pay in partial amounts")
	updateCmd.Flags().String("reference-id", "", "Unique reference ID for internal tracking (max 40 characters)")
	updateCmd.Flags().Int64("expire-by", 0, "Link expiry as Unix timestamp (max 6 months from creation date)")
	updateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
