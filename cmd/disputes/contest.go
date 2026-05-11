package disputes

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var contestCmd = &cobra.Command{
	Use:   "contest <dispute_id>",
	Short: "Contest a dispute (use --action submit to finalise)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		params, _ := cmd.Flags().GetStringArray("param")
		body, err := api.ParseParams(params)
		if err != nil {
			return err
		}
		action, _ := cmd.Flags().GetString("action")
		if action != "" {
			if action != "draft" && action != "submit" {
				return fmt.Errorf("--action must be draft or submit")
			}
			body["action"] = action
		}
		data, err := client.Patch(basePath+"/"+args[0]+"/contest", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(contestCmd)

	contestCmd.Flags().StringArray("param", nil, "Parameter as key=value (e.g. --param amount=1000)")
	contestCmd.Flags().String("action", "", "Contest action: draft (save) or submit (finalise)")
}
