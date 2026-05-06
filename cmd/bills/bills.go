package bills

import "github.com/spf13/cobra"

const basePath = "/v1/bills"

var Cmd = &cobra.Command{
	Use:   "bills",
	Short: "Manage bills",
}
