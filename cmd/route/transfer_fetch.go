package route

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferFetchCmd = &cobra.Command{
	Use:   "fetch <transfer_id>",
	Short: "Fetch a transfer by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if expand, _ := cmd.Flags().GetBool("expand-settlement"); expand {
			q.Set("expand[]", "recipient_settlement")
		}
		data, err := client.Get(transfersPath+"/"+args[0], q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferFetchCmd)
	transferFetchCmd.Flags().Bool("expand-settlement", false, "Include recipient_settlement details in the response")
}
