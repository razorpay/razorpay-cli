package documents

import (
	"fmt"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Upload a document",
	Long: `Upload a document using multipart/form-data.

Example:
  razorpay documents create --file /path/to/file.jpg --purpose dispute_evidence`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()

		filePath, _ := cmd.Flags().GetString("file")
		purpose, _ := cmd.Flags().GetString("purpose")

		if filePath == "" {
			return fmt.Errorf("--file is required")
		}
		if purpose == "" {
			return fmt.Errorf("--purpose is required")
		}

		data, err := client.PostMultipart(basePath, filePath, map[string]string{
			"purpose": purpose,
		})
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(createCmd)

	createCmd.Flags().String("file", "", "Path to the file to upload (required). Supported: jpg, jpeg, png, pdf")
	createCmd.Flags().String("purpose", "", "Purpose of the document e.g. dispute_evidence (required)")
}
