package payments

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <payment_id>",
	Short: "Fetch a payment by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if expands, _ := cmd.Flags().GetStringArray("expand"); len(expands) > 0 {
			for _, e := range expands {
				q.Add("expand[]", e)
			}
		}
		data, err := client.Get(basePath+"/"+args[0], q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(fetchCmd)

	fetchCmd.Flags().StringArray("expand", nil, "Expand related objects (e.g. --expand card --expand emi --expand offers)")
}
