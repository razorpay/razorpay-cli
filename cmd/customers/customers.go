package customers

import "github.com/spf13/cobra"

const basePath = "/v1/customers"

// Cmd is the root customers command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "customers",
	Short: "Manage customers",
}
