package smartcollect

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchByUTRCmd = &cobra.Command{
	Use:   "fetch-by-utr",
	Short: "Fetch payments using a UTR (Unique Transaction Reference) number",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaTransactionID, _ := cmd.Flags().GetString("va-transaction-id")
		if vaTransactionID == "" {
			return fmt.Errorf("--va-transaction-id is required")
		}
		client := cmdutil.NewClient()
		q := url.Values{}
		q.Set("va_transaction_id", vaTransactionID)
		q.Set("virtual_account", "1")
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		data, err := client.Get("/payments", q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(fetchByUTRCmd)

	fetchByUTRCmd.Flags().String("va-transaction-id", "", "Virtual account transaction ID / UTR (required)")
	fetchByUTRCmd.Flags().Int("count", 10, "Number of payments to fetch (max 100)")
	fetchByUTRCmd.Flags().Int("skip", 0, "Number of payments to skip")
}
