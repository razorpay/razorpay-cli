package invoices

import "github.com/spf13/cobra"

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage invoice items",
}

func init() {
	Cmd.AddCommand(itemsCmd)
}
