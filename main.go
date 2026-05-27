package main

import "github.com/razorpay/razorpay-cli/cmd"

// version, commit, and date are set at build time by goreleaser via -ldflags.
// Defaults are used when the binary is built without goreleaser (e.g. go run .).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}
