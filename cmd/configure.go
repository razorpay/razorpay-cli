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

var (
	configureKeyID     string
	configureKeySecret string
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Razorpay API credentials",
	Long: `Configure your Razorpay API Key ID and Key Secret.

Credentials are stored in ~/.razorpay/config.yaml

You can provide credentials via flags:
  razorpay configure --key-id rzp_test_xxxxxxxxxxxx --key-secret xxxxxxxxxxxxxxxxxxxx

Any flag you omit will be prompted for interactively.

You can also set credentials via environment variables:
  RAZORPAY_KEY_ID
  RAZORPAY_KEY_SECRET`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Init()
		reader := bufio.NewReader(os.Stdin)

		keyID := strings.TrimSpace(configureKeyID)
		if keyID == "" {
			input, err := promptValue(reader, "Razorpay Key ID", config.KeyID(), false)
			if err != nil {
				return err
			}
			keyID = input
		}

		keySecret := strings.TrimSpace(configureKeySecret)
		if keySecret == "" {
			input, err := promptValue(reader, "Razorpay Key Secret", config.KeySecret(), true)
			if err != nil {
				return err
			}
			keySecret = input
		}

		if keyID == "" || keySecret == "" {
			return fmt.Errorf("key ID and key secret cannot be empty")
		}

		if err := config.Save(keyID, keySecret); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("\nCredentials saved to %s\n", config.ConfigFilePath())
		return nil
	},
}

// promptValue renders an AWS-style prompt: `Label [hint]: `. If `secret` is
// true and stdin is a TTY, input is read without echoing. When the user
// submits an empty line, the existing value is kept.
func promptValue(reader *bufio.Reader, label, existing string, secret bool) (string, error) {
	fmt.Printf("%s [%s]: ", label, maskedHint(existing, secret))

	var input string
	if secret && term.IsTerminal(int(syscall.Stdin)) {
		b, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return "", err
		}
		input = string(b)
	} else {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = line
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return existing, nil
	}
	return input, nil
}

// maskedHint returns the hint shown inside `[...]` next to a prompt:
//   - "None" when there is no existing value
//   - last 4 characters preceded by `****` for secrets
//   - the value itself for non-secrets
func maskedHint(value string, secret bool) string {
	if value == "" {
		return "None"
	}
	if !secret {
		return value
	}
	if len(value) <= 4 {
		return "****"
	}
	return "****" + value[len(value)-4:]
}

func init() {
	configureCmd.Flags().StringVar(&configureKeyID, "key-id", "", "Razorpay API Key ID")
	configureCmd.Flags().StringVar(&configureKeySecret, "key-secret", "", "Razorpay API Key Secret")
}
