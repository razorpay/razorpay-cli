package subscriptions

import "github.com/spf13/cobra"

const (
	basePath  = "/subscriptions"
	plansPath = "/plans"
)

var Cmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage subscriptions and plans",
}
