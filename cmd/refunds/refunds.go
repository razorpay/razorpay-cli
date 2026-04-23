package refunds

import "github.com/spf13/cobra"

const basePath = "/refunds"

// Cmd is the root refunds command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "refunds",
	Short: "Manage refunds",
}
