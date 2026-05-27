package route

import "github.com/spf13/cobra"

const (
	transfersPath = "/v1/transfers"
	accountsPath  = "/v2/accounts"
)

var Cmd = &cobra.Command{
	Use:   "route",
	Short: "Manage Route transfers and linked accounts",
}
