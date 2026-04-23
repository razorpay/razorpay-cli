package settlements

import "github.com/spf13/cobra"

const (
	basePath     = "/settlements"
	ondemandPath = basePath + "/ondemand"
	reconPath    = basePath + "/recon/combined"
)

// Cmd is the root settlements command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "settlements",
	Short: "Manage settlements",
}
