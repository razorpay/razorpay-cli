package documents

import (
	"fmt"
	"os"

	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchContentCmd = &cobra.Command{
	Use:   "fetch-content <document_id>",
	Short: "Download the content of a document",
	Long: `Download the binary content of a previously uploaded document.

Example:
  razorpay documents fetch-content doc_1234567890abcd --output /tmp/file.jpg`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		outputPath, _ := cmd.Flags().GetString("output")

		data, err := client.Get(basePath+"/"+args[0]+"/content", nil)
		if err != nil {
			cmdutil.HandleErr(err)
			return nil
		}

		if outputPath != "" {
			if err := os.WriteFile(outputPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("Document content saved to %s\n", outputPath)
		} else {
			os.Stdout.Write(data)
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(fetchContentCmd)

	fetchContentCmd.Flags().String("output", "", "Path to save the downloaded file (prints to stdout if not set)")
}
