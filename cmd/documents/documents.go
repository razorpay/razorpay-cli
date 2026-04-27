package documents

import "github.com/spf13/cobra"

const basePath = "/v1/documents"

// Cmd is the root documents command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "documents",
	Short: "Manage documents",
}
