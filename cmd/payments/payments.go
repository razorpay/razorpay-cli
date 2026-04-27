package payments

import "github.com/spf13/cobra"

const basePath = "/v1/payments"

// Cmd is the root payments command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage payments",
}
