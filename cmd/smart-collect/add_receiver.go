package smartcollect

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var addReceiverCmd = &cobra.Command{
	Use:   "add-receiver <virtual_account_id>",
	Short: "Add a receiver (bank_account or vpa) to an existing virtual account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		receiverType, _ := cmd.Flags().GetString("type")
		if receiverType != "bank_account" && receiverType != "vpa" {
			return fmt.Errorf("--type must be bank_account or vpa")
		}
		vpaDescriptor, _ := cmd.Flags().GetString("vpa-descriptor")
		client := cmdutil.NewClient()
		body := map[string]any{
			"types": []string{receiverType},
		}
		if receiverType == "vpa" && vpaDescriptor != "" {
			body["vpa"] = map[string]any{
				"descriptor": vpaDescriptor,
			}
		}
		data, err := client.Post(basePath+"/"+args[0]+"/receivers", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(addReceiverCmd)

	addReceiverCmd.Flags().String("type", "", "Receiver type to add: bank_account or vpa (required)")
	_ = addReceiverCmd.MarkFlagRequired("type")
	addReceiverCmd.Flags().String("vpa-descriptor", "", "VPA descriptor / label (used when --type=vpa)")
}
