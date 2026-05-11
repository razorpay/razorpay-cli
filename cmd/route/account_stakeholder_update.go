package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var stakeholderUpdateCmd = &cobra.Command{
	Use:   "stakeholder-update <account_id> <stakeholder_id>",
	Short: "Update a stakeholder for a linked account",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("email"); v != "" {
			body["email"] = v
		}
		if v, _ := cmd.Flags().GetFloat64("percentage-ownership"); v > 0 {
			body["percentage_ownership"] = v
		}

		relationship := map[string]any{}
		if cmd.Flags().Changed("is-director") {
			v, _ := cmd.Flags().GetBool("is-director")
			relationship["director"] = v
		}
		if cmd.Flags().Changed("is-executive") {
			v, _ := cmd.Flags().GetBool("is-executive")
			relationship["executive"] = v
		}
		if len(relationship) > 0 {
			body["relationship"] = relationship
		}

		phone := map[string]any{}
		if v, _ := cmd.Flags().GetString("phone-primary"); v != "" {
			phone["primary"] = v
		}
		if v, _ := cmd.Flags().GetString("phone-secondary"); v != "" {
			phone["secondary"] = v
		}
		if len(phone) > 0 {
			body["phone"] = phone
		}

		addr := map[string]any{}
		for src, dst := range map[string]string{
			"street":      "street",
			"city":        "city",
			"state":       "state",
			"postal-code": "postal_code",
			"country":     "country",
		} {
			if v, _ := cmd.Flags().GetString(src); v != "" {
				addr[dst] = v
			}
		}
		if len(addr) > 0 {
			body["addresses"] = map[string]any{"residential": addr}
		}

		if pan, _ := cmd.Flags().GetString("pan"); pan != "" {
			body["kyc"] = map[string]any{"pan": pan}
		}

		if notes, _ := cmd.Flags().GetStringArray("note"); len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag must be provided to update")
		}

		data, err := client.Patch(accountsPath+"/"+args[0]+"/stakeholders/"+args[1], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(stakeholderUpdateCmd)

	stakeholderUpdateCmd.Flags().String("name", "", "Stakeholder name as per PAN card")
	stakeholderUpdateCmd.Flags().String("email", "", "Stakeholder email address")
	stakeholderUpdateCmd.Flags().Float64("percentage-ownership", 0, "Ownership percentage (max 2 decimal places)")
	stakeholderUpdateCmd.Flags().Bool("is-director", false, "Stakeholder is a director")
	stakeholderUpdateCmd.Flags().Bool("is-executive", false, "Stakeholder is an executive")
	stakeholderUpdateCmd.Flags().String("phone-primary", "", "Primary phone number")
	stakeholderUpdateCmd.Flags().String("phone-secondary", "", "Secondary phone number")
	stakeholderUpdateCmd.Flags().String("street", "", "Residential street address")
	stakeholderUpdateCmd.Flags().String("city", "", "Residential city")
	stakeholderUpdateCmd.Flags().String("state", "", "Residential state")
	stakeholderUpdateCmd.Flags().String("postal-code", "", "Residential postal code")
	stakeholderUpdateCmd.Flags().String("country", "", "Residential country code")
	stakeholderUpdateCmd.Flags().String("pan", "", "PAN number for KYC")
	stakeholderUpdateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
