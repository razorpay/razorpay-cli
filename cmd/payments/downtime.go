package payments

import "github.com/spf13/cobra"

var downtimeCmd = &cobra.Command{
	Use:   "downtime",
	Short: "Manage payment method downtimes",
}

func init() {
	Cmd.AddCommand(downtimeCmd)
}
