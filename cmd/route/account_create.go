package route

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a linked account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		email, _ := cmd.Flags().GetString("email")
		phone, _ := cmd.Flags().GetString("phone")
		legalBusinessName, _ := cmd.Flags().GetString("legal-business-name")
		businessType, _ := cmd.Flags().GetString("business-type")
		contactName, _ := cmd.Flags().GetString("contact-name")
		customerFacingName, _ := cmd.Flags().GetString("customer-facing-business-name")
		referenceID, _ := cmd.Flags().GetString("reference-id")

		if email == "" {
			return fmt.Errorf("--email is required")
		}
		if phone == "" {
			return fmt.Errorf("--phone is required")
		}
		if legalBusinessName == "" {
			return fmt.Errorf("--legal-business-name is required")
		}
		if businessType == "" {
			return fmt.Errorf("--business-type is required")
		}
		if contactName == "" {
			return fmt.Errorf("--contact-name is required")
		}

		body := map[string]any{
			"email":               email,
			"phone":               phone,
			"legal_business_name": legalBusinessName,
			"business_type":       businessType,
			"contact_name":        contactName,
		}
		if customerFacingName != "" {
			body["customer_facing_business_name"] = customerFacingName
		}
		if referenceID != "" {
			body["reference_id"] = referenceID
		}

		// profile object
		profileCategory, _ := cmd.Flags().GetString("profile-category")
		profileSubcategory, _ := cmd.Flags().GetString("profile-subcategory")
		profileBusinessModel, _ := cmd.Flags().GetString("profile-business-model")
		profile := map[string]any{}
		if profileCategory != "" {
			profile["category"] = profileCategory
		}
		if profileSubcategory != "" {
			profile["subcategory"] = profileSubcategory
		}
		if profileBusinessModel != "" {
			profile["business_model"] = profileBusinessModel
		}
		// registered address
		regStreet1, _ := cmd.Flags().GetString("registered-street1")
		regStreet2, _ := cmd.Flags().GetString("registered-street2")
		regCity, _ := cmd.Flags().GetString("registered-city")
		regState, _ := cmd.Flags().GetString("registered-state")
		regPostal, _ := cmd.Flags().GetString("registered-postal-code")
		regCountry, _ := cmd.Flags().GetString("registered-country")
		regAddr := map[string]any{}
		if regStreet1 != "" {
			regAddr["street1"] = regStreet1
		}
		if regStreet2 != "" {
			regAddr["street2"] = regStreet2
		}
		if regCity != "" {
			regAddr["city"] = regCity
		}
		if regState != "" {
			regAddr["state"] = regState
		}
		if regPostal != "" {
			regAddr["postal_code"] = regPostal
		}
		if regCountry != "" {
			regAddr["country"] = regCountry
		}
		if len(regAddr) > 0 {
			profile["addresses"] = map[string]any{"registered": regAddr}
		}
		if len(profile) > 0 {
			body["profile"] = profile
		}

		// legal_info
		pan, _ := cmd.Flags().GetString("pan")
		gst, _ := cmd.Flags().GetString("gst")
		legalInfo := map[string]any{}
		if pan != "" {
			legalInfo["pan"] = pan
		}
		if gst != "" {
			legalInfo["gst"] = gst
		}
		if len(legalInfo) > 0 {
			body["legal_info"] = legalInfo
		}

		// contact_info
		contactInfo := map[string]any{}
		buildContactType := func(typeKey, emailFlag, phoneFlag, policyFlag string) {
			e, _ := cmd.Flags().GetString(emailFlag)
			p, _ := cmd.Flags().GetString(phoneFlag)
			u, _ := cmd.Flags().GetString(policyFlag)
			m := map[string]any{}
			if e != "" {
				m["email"] = e
			}
			if p != "" {
				m["phone"] = p
			}
			if u != "" {
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
		websites, _ := cmd.Flags().GetStringArray("website")
		if len(websites) > 0 {
			body["apps"] = map[string]any{"websites": websites}
		}

		// notes
		notes, _ := cmd.Flags().GetStringArray("note")
		if len(notes) > 0 {
			notesMap, err := api.ParseParams(notes)
			if err != nil {
				return err
			}
			body["notes"] = notesMap
		}

		data, err := client.Post(accountsPath, body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(accountCreateCmd)

	// Required
	accountCreateCmd.Flags().String("email", "", "Linked account business email (required)")
	accountCreateCmd.Flags().String("phone", "", "Business phone number (required)")
	accountCreateCmd.Flags().String("legal-business-name", "", "Registered business name (required)")
	accountCreateCmd.Flags().String("business-type", "", "Business type, e.g. route, individual, sole_proprietorship, partnership, private_limited (required)")
	accountCreateCmd.Flags().String("contact-name", "", "Contact person name (required)")

	// Optional top-level
	accountCreateCmd.Flags().String("customer-facing-business-name", "", "Public-facing business name shown on checkout")
	accountCreateCmd.Flags().String("reference-id", "", "Your internal reference ID for the account")

	// profile
	accountCreateCmd.Flags().String("profile-category", "", "Business category")
	accountCreateCmd.Flags().String("profile-subcategory", "", "Business subcategory")
	accountCreateCmd.Flags().String("profile-business-model", "", "Description of the business model")
	accountCreateCmd.Flags().String("registered-street1", "", "Registered address line 1")
	accountCreateCmd.Flags().String("registered-street2", "", "Registered address line 2")
	accountCreateCmd.Flags().String("registered-city", "", "Registered address city")
	accountCreateCmd.Flags().String("registered-state", "", "Registered address state")
	accountCreateCmd.Flags().String("registered-postal-code", "", "Registered address postal code")
	accountCreateCmd.Flags().String("registered-country", "", "Registered address country code")

	// legal_info
	accountCreateCmd.Flags().String("pan", "", "Business PAN number")
	accountCreateCmd.Flags().String("gst", "", "Business GST number")

	// contact_info
	accountCreateCmd.Flags().String("chargeback-email", "", "Chargeback contact email")
	accountCreateCmd.Flags().String("chargeback-phone", "", "Chargeback contact phone")
	accountCreateCmd.Flags().String("chargeback-policy-url", "", "Chargeback policy URL")
	accountCreateCmd.Flags().String("refund-email", "", "Refund contact email")
	accountCreateCmd.Flags().String("refund-phone", "", "Refund contact phone")
	accountCreateCmd.Flags().String("refund-policy-url", "", "Refund policy URL")
	accountCreateCmd.Flags().String("support-email", "", "Support contact email")
	accountCreateCmd.Flags().String("support-phone", "", "Support contact phone")
	accountCreateCmd.Flags().String("support-policy-url", "", "Support policy URL")

	// apps
	accountCreateCmd.Flags().StringArray("website", nil, "Website URL (repeatable)")

	// notes
	accountCreateCmd.Flags().StringArray("note", nil, "Note as key=value (repeatable, max 15 pairs)")
}
