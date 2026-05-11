package subscriptions

import "github.com/spf13/cobra"

const (
	basePath  = "/v1/subscriptions"
	plansPath = "/v1/plans"
)

var Cmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage subscriptions and plans",
}
