package subscriptions

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume <subscription_id>",
	Short: "Resume a paused subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resumeAt, _ := cmd.Flags().GetString("resume-at")
		if resumeAt != "now" && resumeAt != "cycle_end" {
			return fmt.Errorf("--resume-at must be now or cycle_end")
		}
		client := cmdutil.NewClient()
		body := map[string]interface{}{
			"resume_at": resumeAt,
		}
		data, err := client.Post(basePath+"/"+args[0]+"/resume", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(resumeCmd)

	resumeCmd.Flags().String("resume-at", "now", "When to resume: now")
}
