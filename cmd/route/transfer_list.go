package route

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var transferListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transfers",
	Long: `List all transfers. Optionally filter by settlement ID to fetch
	transfers for a specific settlement, or use --expand to include
	settlement details in the response.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if expands, _ := cmd.Flags().GetStringArray("expand"); len(expands) > 0 {
			for _, e := range expands {
				q.Add("expand[]", e)
			}
		}
		if transferType, _ := cmd.Flags().GetString("transfer-type"); transferType != "" {
			q.Set("transfer_type", transferType)
		}
		if settlementID, _ := cmd.Flags().GetString("settlement-id"); settlementID != "" {
			q.Set("recipient_settlement_id", settlementID)
		}
		data, err := client.Get(transfersPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(transferListCmd)

	transferListCmd.Flags().StringArray("expand", nil, "Expand related objects (e.g. --expand recipient_settlement)")
	transferListCmd.Flags().String("transfer-type", "", "Transfer type filter for partners: platform or regular")
	transferListCmd.Flags().String("settlement-id", "", "Filter transfers by settlement ID (from settlement.processed webhook)")
}
