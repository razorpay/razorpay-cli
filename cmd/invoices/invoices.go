package invoices

import "github.com/spf13/cobra"

const (
	basePath  = "/v1/invoices"
	itemsPath = "/v1/items"
)

var Cmd = &cobra.Command{
	Use:   "invoices",
	Short: "Manage invoices",
}

// buildAddress reads --<prefix>-line1/line2/zipcode/city/state/country flags
// and returns a populated address map (or nil if no flags were set).
func buildAddress(cmd *cobra.Command, prefix string) map[string]any {
	addr := map[string]any{}
	for _, field := range []string{"line1", "line2", "zipcode", "city", "state", "country"} {
		if v, _ := cmd.Flags().GetString(prefix + "-" + field); v != "" {
			addr[field] = v
		}
	}
	return addr
}
