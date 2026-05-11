package invoices

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify <invoice_id>",
	Short: "Send a notification for an invoice (via sms or email)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		medium, _ := cmd.Flags().GetString("medium")
		if medium != "sms" && medium != "email" {
			return fmt.Errorf("--medium must be sms or email")
		}
		client := cmdutil.NewClient()
		data, err := client.Post(basePath+"/"+args[0]+"/notify_by/"+medium, nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(notifyCmd)

	notifyCmd.Flags().String("medium", "", "Notification medium: sms or email (required)")
	_ = notifyCmd.MarkFlagRequired("medium")
}
