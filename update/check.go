package update

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const checkInterval = 24 * time.Hour

// CheckOnce prints a notice to stderr if a newer release exists. Throttled to
// once per checkInterval via ~/.razorpay/version-check.json. Silent on every
// failure path.
func CheckOnce(currentVersion string) {
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		return
	}

	c := readCache()
	// elapsed > 0 guards against a future-dated cache (clock skew) trapping us.
	elapsed := time.Since(c.CheckedAt)
	if elapsed > 0 && elapsed < checkInterval {
		return
	}

	latest, err := fetchLatestVersion(context.Background())
	if err != nil {
		return
	}

	_ = writeCache(cache{CheckedAt: time.Now().UTC()})

	if isNewer(latest, currentVersion) {
		fmt.Fprintf(os.Stderr,
			"\nA new version of razorpay is available: %s. See https://razorpay.com/docs/api/install-cli/ to upgrade.\n\n",
			latest)
	}
}

// isNewer returns true when current is "dev" (unstamped build) or differs
// from latest after trimming the optional "v" prefix.
func isNewer(latest, current string) bool {
	if latest == "" {
		return false
	}
	if current == "" || current == "dev" {
		return true
	}
	return strings.TrimPrefix(latest, "v") != strings.TrimPrefix(current, "v")
}
