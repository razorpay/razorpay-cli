package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var accountUpdateCmd = &cobra.Command{
	Use:   "update <account_id>",
	Short: "Update a linked account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		if v, _ := cmd.Flags().GetString("phone"); v != "" {
			body["phone"] = v
		}
		if v, _ := cmd.Flags().GetString("customer-facing-business-name"); v != "" {
			body["customer_facing_business_name"] = v
		}

		// profile object
		profile := map[string]any{}
		if v, _ := cmd.Flags().GetString("profile-category"); v != "" {
			profile["category"] = v
		}
		if v, _ := cmd.Flags().GetString("profile-subcategory"); v != "" {
			profile["subcategory"] = v
		}
		if v, _ := cmd.Flags().GetString("profile-business-model"); v != "" {
			profile["business_model"] = v
		}
		regAddr := map[string]any{}
		for src, dst := range map[string]string{
			"registered-street1":      "street1",
			"registered-street2":      "street2",
			"registered-city":         "city",
			"registered-state":        "state",
			"registered-postal-code":  "postal_code",
			"registered-country":      "country",
		} {
			if v, _ := cmd.Flags().GetString(src); v != "" {
				regAddr[dst] = v
			}
		}
		if len(regAddr) > 0 {
			profile["addresses"] = map[string]any{"registered": regAddr}
		}
		if len(profile) > 0 {
			body["profile"] = profile
		}

		// legal_info
		legalInfo := map[string]any{}
		if v, _ := cmd.Flags().GetString("pan"); v != "" {
			legalInfo["pan"] = v
		}
		if v, _ := cmd.Flags().GetString("gst"); v != "" {
			legalInfo["gst"] = v
		}
		if len(legalInfo) > 0 {
			body["legal_info"] = legalInfo
		}

		// contact_info
		contactInfo := map[string]any{}
		buildContactType := func(typeKey, emailFlag, phoneFlag, policyFlag string) {
			m := map[string]any{}
			if e, _ := cmd.Flags().GetString(emailFlag); e != "" {
				m["email"] = e
			}
			if p, _ := cmd.Flags().GetString(phoneFlag); p != "" {
				m["phone"] = p
			}
			if u, _ := cmd.Flags().GetString(policyFlag); u != "" {
				m["policy_url"] = u
			}
			if len(m) > 0 {
				contactInfo[typeKey] = m
			}
		}
		buildContactType("chargeback", "chargeback-email", "chargeback-phone", "chargeback-policy-url")
		buildContactType("refund", "refund-email", "refund-phone", "refund-policy-url")
		buildContactType("support", "support-email", "support-phone", "support-policy-url")
		if len(contactInfo) > 0 {
			body["contact_info"] = contactInfo
		}

		// apps
		if websites, _ := cmd.Flags().GetStringArray("website"); len(websites) > 0 {
			body["apps"] = map[string]any{"websites": websites}
		}

		// notes
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

		data, err := client.Patch(accountsPath+"/"+args[0], body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(accountUpdateCmd)

	accountUpdateCmd.Flags().String("phone", "", "Business phone number (8-15 digits)")
	accountUpdateCmd.Flags().String("customer-facing-business-name", "", "Public-facing business name shown on checkout")
	accountUpdateCmd.Flags().String("profile-category", "", "Business category")
	accountUpdateCmd.Flags().String("profile-subcategory", "", "Business subcategory")
	accountUpdateCmd.Flags().String("profile-business-model", "", "Description of the business model")
	accountUpdateCmd.Flags().String("registered-street1", "", "Registered address line 1")
	accountUpdateCmd.Flags().String("registered-street2", "", "Registered address line 2")
	accountUpdateCmd.Flags().String("registered-city", "", "Registered address city")
	accountUpdateCmd.Flags().String("registered-state", "", "Registered address state")
	accountUpdateCmd.Flags().String("registered-postal-code", "", "Registered address postal code")
	accountUpdateCmd.Flags().String("registered-country", "", "Registered address country code")
	accountUpdateCmd.Flags().String("pan", "", "Business PAN number")
	accountUpdateCmd.Flags().String("gst", "", "Business GST number")
	accountUpdateCmd.Flags().String("chargeback-email", "", "Chargeback contact email")
	accountUpdateCmd.Flags().String("chargeback-phone", "", "Chargeback contact phone")
	accountUpdateCmd.Flags().String("chargeback-policy-url", "", "Chargeback policy URL")
	accountUpdateCmd.Flags().String("refund-email", "", "Refund contact email")
	accountUpdateCmd.Flags().String("refund-phone", "", "Refund contact phone")
	accountUpdateCmd.Flags().String("refund-policy-url", "", "Refund policy URL")
	accountUpdateCmd.Flags().String("support-email", "", "Support contact email")
	accountUpdateCmd.Flags().String("support-phone", "", "Support contact phone")
	accountUpdateCmd.Flags().String("support-policy-url", "", "Support policy URL")
	accountUpdateCmd.Flags().StringArray("website", nil, "Website URL (repeatable)")
	accountUpdateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
