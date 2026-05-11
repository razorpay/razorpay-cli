package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var stakeholderCreateCmd = &cobra.Command{
	Use:   "stakeholder-create <account_id>",
	Short: "Create a stakeholder for a linked account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		email, _ := cmd.Flags().GetString("email")
		name, _ := cmd.Flags().GetString("name")
		if email == "" {
			return fmt.Errorf("--email is required")
		}
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		body := map[string]any{
			"name":  name,
			"email": email,
		}

		if v, _ := cmd.Flags().GetFloat64("percentage-ownership"); v > 0 {
			body["percentage_ownership"] = v
		}

		// relationship
		relationship := map[string]any{}
		if cmd.Flags().Changed("relationship-director") {
			v, _ := cmd.Flags().GetBool("relationship-director")
			relationship["director"] = v
		}
		if cmd.Flags().Changed("relationship-executive") {
			v, _ := cmd.Flags().GetBool("relationship-executive")
			relationship["executive"] = v
		}
		if len(relationship) > 0 {
			body["relationship"] = relationship
		}

		// phone
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

		// residential address
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

		// KYC
		if pan, _ := cmd.Flags().GetString("pan"); pan != "" {
			body["kyc"] = map[string]any{"pan": pan}
		}

		// notes
		if notes, _ := cmd.Flags().GetStringArray("note"); len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		data, err := client.Post(accountsPath+"/"+args[0]+"/stakeholders", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(stakeholderCreateCmd)

	stakeholderCreateCmd.Flags().String("name", "", "Stakeholder name as per PAN card (required)")
	stakeholderCreateCmd.Flags().String("email", "", "Stakeholder email address (required)")
	stakeholderCreateCmd.Flags().Float64("percentage-ownership", 0, "Ownership percentage (max 2 decimal places)")
	stakeholderCreateCmd.Flags().Bool("relationship-director", false, "Stakeholder is a director")
	stakeholderCreateCmd.Flags().Bool("relationship-executive", false, "Stakeholder is an executive")
	stakeholderCreateCmd.Flags().String("phone-primary", "", "Primary phone number")
	stakeholderCreateCmd.Flags().String("phone-secondary", "", "Secondary phone number")
	stakeholderCreateCmd.Flags().String("street", "", "Residential street address")
	stakeholderCreateCmd.Flags().String("city", "", "Residential city")
	stakeholderCreateCmd.Flags().String("state", "", "Residential state")
	stakeholderCreateCmd.Flags().String("postal-code", "", "Residential postal code")
	stakeholderCreateCmd.Flags().String("country", "", "Residential country code")
	stakeholderCreateCmd.Flags().String("pan", "", "PAN number for KYC")
	stakeholderCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
