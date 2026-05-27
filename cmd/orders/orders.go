package orders

import "github.com/spf13/cobra"

const basePath = "/v1/orders"

// Cmd is the root orders command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage orders",
}
