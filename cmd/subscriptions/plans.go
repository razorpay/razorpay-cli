package subscriptions

import "github.com/spf13/cobra"

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Manage subscription plans",
}

func init() {
	Cmd.AddCommand(plansCmd)
}
