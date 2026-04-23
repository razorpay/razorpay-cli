package smartcollect

import "github.com/spf13/cobra"

const basePath = "/virtual_accounts"

var Cmd = &cobra.Command{
	Use:   "smart-collect",
	Short: "Manage Smart Collect virtual accounts",
}
