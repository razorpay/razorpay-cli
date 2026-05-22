package disputes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var contestCmd = &cobra.Command{
	Use:   "contest <dispute_id>",
	Short: "Contest a dispute with evidence documents",
	Long: `Contest a dispute by providing evidence documents and an optional summary.

Evidence document flags accept document IDs (repeatable for multiple docs).
Use --action submit to finalise, or --action draft (default) to save progress.

For custom evidence types, use --others as a JSON array:
  --others '[{"type":"receipt_signed_by_customer","document_ids":["doc_abc","doc_def"]}]'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		body := map[string]any{}

		if amount, _ := cmd.Flags().GetInt("amount"); cmd.Flags().Changed("amount") {
			body["amount"] = amount
		}
		if summary, _ := cmd.Flags().GetString("summary"); summary != "" {
			body["summary"] = summary
		}

		// Evidence document fields (all are lists of document IDs)
		evidenceFields := []string{
			"shipping_proof",
			"billing_proof",
			"cancellation_proof",
			"customer_communication",
			"proof_of_service",
			"explanation_letter",
			"refund_confirmation",
			"access_activity_log",
			"refund_cancellation_policy",
			"term_and_conditions",
		}
		for _, field := range evidenceFields {
			flag := strings.ReplaceAll(field, "_", "-")
			if docs, _ := cmd.Flags().GetStringArray(flag); len(docs) > 0 {
				body[field] = docs
			}
		}

		// Others: JSON array of {type, document_ids} objects
		if othersJSON, _ := cmd.Flags().GetString("others"); othersJSON != "" {
			var others any
			if err := json.Unmarshal([]byte(othersJSON), &others); err != nil {
				return fmt.Errorf("--others is not valid JSON: %w", err)
			}
			body["others"] = others
		}

		action, _ := cmd.Flags().GetString("action")
		if action != "" {
			if action != "draft" && action != "submit" {
				return fmt.Errorf("--action must be draft or submit")
			}
			body["action"] = action
		}

		if len(body) == 0 {
			return fmt.Errorf("at least one flag is required (see --help for available evidence types)")
		}

		data, err := client.Patch(basePath+"/"+args[0]+"/contest", body)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(contestCmd)

	contestCmd.Flags().Int("amount", 0, "Contest amount in paise (defaults to full dispute amount if omitted)")
	contestCmd.Flags().String("summary", "", "Explanation for contesting the dispute (max 1000 chars)")
	contestCmd.Flags().StringArray("shipping-proof", nil, "Document ID for shipping proof (repeatable)")
	contestCmd.Flags().StringArray("billing-proof", nil, "Document ID for billing/order confirmation proof (repeatable)")
	contestCmd.Flags().StringArray("cancellation-proof", nil, "Document ID for cancellation proof (repeatable)")
	contestCmd.Flags().StringArray("customer-communication", nil, "Document ID for customer communication proof (repeatable)")
	contestCmd.Flags().StringArray("proof-of-service", nil, "Document ID for proof of service (repeatable)")
	contestCmd.Flags().StringArray("explanation-letter", nil, "Document ID for explanation letter (repeatable)")
	contestCmd.Flags().StringArray("refund-confirmation", nil, "Document ID for refund confirmation proof (repeatable)")
	contestCmd.Flags().StringArray("access-activity-log", nil, "Document ID for access/activity log (repeatable)")
	contestCmd.Flags().StringArray("refund-cancellation-policy", nil, "Document ID for refund/cancellation policy (repeatable)")
	contestCmd.Flags().StringArray("term-and-conditions", nil, "Document ID for terms and conditions (repeatable)")
	contestCmd.Flags().String("others", "", `Custom evidence as JSON array. Example: '[{"type":"receipt_signed_by_customer","document_ids":["doc_abc","doc_def"]}]'`)
	contestCmd.Flags().String("action", "", "Contest action: draft (save) or submit (finalise)")
}
