package smartcollect

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <virtual_account_id>",
	Short: "Update a virtual account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		if v, _ := cmd.Flags().GetInt64("close-by"); v > 0 {
			body["close_by"] = v
		}
		if v, _ := cmd.Flags().GetString("description"); v != "" {
			body["description"] = v
		}
		if notes, _ := cmd.Flags().GetStringArray("note"); len(notes) > 0 {
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

	updateCmd.Flags().Int64("close-by", 0, "New Unix timestamp for automatic closure (min 15 minutes from now)")
	updateCmd.Flags().String("description", "", "Updated description for the virtual account")
	updateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
