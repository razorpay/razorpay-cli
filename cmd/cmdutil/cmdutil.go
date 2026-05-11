package cmdutil

import (
	"fmt"
	"os"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/config"
)

// NewClient initialises config and returns an authenticated API client.
func NewClient() *api.Client {
	config.Init()
	return api.New(config.KeyID(), config.KeySecret())
}

// HandleErr prints the error to stderr and exits with code 1.
func HandleErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
