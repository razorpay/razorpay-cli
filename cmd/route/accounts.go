package route

import "github.com/spf13/cobra"

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage linked accounts",
}

func init() {
	Cmd.AddCommand(accountsCmd)
}
