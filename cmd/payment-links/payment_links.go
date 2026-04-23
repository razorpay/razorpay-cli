package paymentlinks

import "github.com/spf13/cobra"

const basePath = "/payment_links"

var Cmd = &cobra.Command{
	Use:   "payment-links",
	Short: "Manage payment links",
}
