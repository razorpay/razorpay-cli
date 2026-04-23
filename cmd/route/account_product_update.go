package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var productUpdateCmd = &cobra.Command{
	Use:   "product-update <account_id> <product_id>",
	Short: "Update Route product configuration for a linked account",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		settlements := map[string]any{}
		if v, _ := cmd.Flags().GetString("account-number"); v != "" {
			settlements["account_number"] = v
		}
		if v, _ := cmd.Flags().GetString("ifsc"); v != "" {
			settlements["ifsc_code"] = v
		}
		if v, _ := cmd.Flags().GetString("beneficiary-name"); v != "" {
			settlements["beneficiary_name"] = v
		}
		if len(settlements) > 0 {
			body["settlements"] = settlements
		}

		if cmd.Flags().Changed("tnc-accepted") {
			accepted, _ := cmd.Flags().GetBool("tnc-accepted")
			body["tnc_accepted"] = accepted
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
		}

		data, err := client.Patch(accountsPath+"/"+args[0]+"/products/"+args[1], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(productUpdateCmd)

	productUpdateCmd.Flags().String("account-number", "", "Bank account number for settlement")
	productUpdateCmd.Flags().String("ifsc", "", "Bank IFSC code for settlement")
	productUpdateCmd.Flags().String("beneficiary-name", "", "Bank account beneficiary name")
	productUpdateCmd.Flags().Bool("tnc-accepted", false, "Accept the Route product terms and conditions")
}
