package route

import "github.com/spf13/cobra"

var transfersCmd = &cobra.Command{
	Use:   "transfers",
	Short: "Manage Route transfers",
}

func init() {
	Cmd.AddCommand(transfersCmd)
}
