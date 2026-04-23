package disputes

import "github.com/spf13/cobra"

const basePath = "/disputes"

// Cmd is the root disputes command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "disputes",
	Short: "Manage disputes",
}
