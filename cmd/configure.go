package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/razorpay/razorpay-cli/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Razorpay API credentials",
	Long: `Configure your Razorpay API Key ID and Key Secret.

Credentials are stored in ~/.razorpay/config.yaml

You can also set credentials via environment variables:
  RAZORPAY_KEY_ID
  RAZORPAY_KEY_SECRET

Non-interactive (agent/script) usage:
  razorpay configure --key-id=rzp_test_xxx --key-secret=yyy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyID, _ := cmd.Flags().GetString("key-id")
		keySecret, _ := cmd.Flags().GetString("key-secret")

		// Non-interactive mode: both flags provided
		if keyID != "" && keySecret != "" {
			if err := config.Save(keyID, keySecret); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("Credentials saved to %s\n", config.ConfigFilePath())
			return nil
		}

		// If only one flag is provided, error early with a clear message
		if keyID != "" || keySecret != "" {
			return fmt.Errorf("both --key-id and --key-secret must be provided for non-interactive mode")
		}

		// Interactive mode
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Key ID: ")
		var err error
		keyID, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		keyID = strings.TrimSpace(keyID)

		fmt.Print("Key Secret: ")
		if term.IsTerminal(int(syscall.Stdin)) {
			b, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				return err
			}
			keySecret = string(b)
		} else {
			keySecret, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			keySecret = strings.TrimSpace(keySecret)
		}

		if keyID == "" || keySecret == "" {
			return fmt.Errorf("key ID and key secret cannot be empty")
		}

		if err := config.Save(keyID, keySecret); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Credentials saved to %s\n", config.ConfigFilePath())
		return nil
	},
}

func init() {
	configureCmd.Flags().String("key-id", "", "Razorpay Key ID (non-interactive mode)")
	configureCmd.Flags().String("key-secret", "", "Razorpay Key Secret (non-interactive mode)")
}
