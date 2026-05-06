package bills

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const basePath = "/v1/bills"

// Cmd is the root bills command registered by the parent cmd package.
var Cmd = &cobra.Command{
	Use:   "bills",
	Short: "Manage bills (Razorpay Billme)",
}

// parseJSONFlag reads a string flag and parses its value as JSON.
// Empty values return (nil, nil) so callers can omit the field from the body.
// Used for nested object/array body params (customer, line_items, etc.) where
// the doc surface is too deep to expose as flat flags.
func parseJSONFlag(cmd *cobra.Command, name string) (interface{}, error) {
	raw, _ := cmd.Flags().GetString(name)
	if raw == "" {
		return nil, nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, fmt.Errorf("--%s: invalid JSON: %w", name, err)
	}
	return v, nil
}
