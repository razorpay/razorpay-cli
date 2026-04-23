package route

import "github.com/spf13/cobra"

const (
	transfersPath = "/transfers"
	accountsPath  = "/v2/accounts" // Route linked accounts use the v2 API
)

var Cmd = &cobra.Command{
	Use:   "route",
	Short: "Manage Route transfers and linked accounts",
}
