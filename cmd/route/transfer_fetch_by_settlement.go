package route

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchTransfersBySettlementCmd = &cobra.Command{
	Use:   "fetch-by-settlement",
	Short: "Fetch all transfers for a specific settlement",
	RunE: func(cmd *cobra.Command, args []string) error {
		settlementID, _ := cmd.Flags().GetString("settlement-id")
		if settlementID == "" {
			return fmt.Errorf("--settlement-id is required")
		}
		client := cmdutil.NewClient()
		q := url.Values{}
		q.Set("recipient_settlement_id", settlementID)
		data, err := client.Get(transfersPath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(fetchTransfersBySettlementCmd)

	fetchTransfersBySettlementCmd.Flags().String("settlement-id", "", "Settlement ID from the settlement.processed webhook (required)")
}
