package disputes

import "github.com/spf13/cobra"

const basePath = "/v1/disputes"

// Cmd is the root disputes command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "disputes",
	Short: "Manage disputes",
}
