package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var productRequestCmd = &cobra.Command{
	Use:   "product-request <account_id>",
	Short: "Request Route product configuration for a linked account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		body := map[string]any{
			"product_name": "route",
		}

		if cmd.Flags().Changed("tnc-accepted") {
			accepted, _ := cmd.Flags().GetBool("tnc-accepted")
			body["tnc_accepted"] = accepted
		}

		data, err := client.Post(accountsPath+"/"+args[0]+"/products", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(productRequestCmd)

	productRequestCmd.Flags().Bool("tnc-accepted", false, "Accept the Route product terms and conditions")
}
