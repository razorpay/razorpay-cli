//go:build e2e
// +build e2e

// Package tests runs the Razorpay CLI as a subprocess against the real
// Razorpay test API and asserts on the responses. It is gated by the `e2e`
// build tag so a plain `go test ./...` never invokes it.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binPath   string
	keyID     string
	keySecret string
)

func TestMain(m *testing.M) {
	keyID = firstNonEmpty(os.Getenv("RAZORPAY_TEST_KEY_ID"), os.Getenv("RAZORPAY_KEY_ID"))
	keySecret = firstNonEmpty(os.Getenv("RAZORPAY_TEST_KEY_SECRET"), os.Getenv("RAZORPAY_KEY_SECRET"))

	if keyID == "" || keySecret == "" {
		fmt.Fprintln(os.Stderr,
			"e2e: RAZORPAY_TEST_KEY_ID and RAZORPAY_TEST_KEY_SECRET are required "+
				"(RAZORPAY_KEY_ID / RAZORPAY_KEY_SECRET are accepted as a fallback). Skipping.")
		os.Exit(0)
	}
	if !strings.HasPrefix(keyID, "rzp_test_") {
		fmt.Fprintln(os.Stderr,
			"e2e: refusing to run against a non-test key; key id must start with 'rzp_test_'.")
		os.Exit(1)
	}

	tmpBin, err := os.MkdirTemp("", "razorpay-cli-e2e-bin-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e: cannot create temp bin dir:", err)
		os.Exit(1)
	}
	binPath = filepath.Join(tmpBin, "razorpay")

	_, thisFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(thisFile))

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = projectRoot
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "e2e: failed to build CLI binary:", err)
		os.Exit(1)
	}

	code := m.Run()
	_ = os.RemoveAll(tmpBin)
	os.Exit(code)
}

type result struct {
	stdout string
	stderr string
	err    error
}

// run executes the CLI with credentials passed via env. Each call gets a fresh
// temp HOME so configure-writes from one test do not leak into another.
func run(t *testing.T, args ...string) result {
	t.Helper()
	return runWithStdin(t, nil, args...)
}

func runWithStdin(t *testing.T, stdin io.Reader, args ...string) result {
	t.Helper()

	tmpHome, err := os.MkdirTemp("", "razorpay-cli-home-")
	if err != nil {
		t.Fatalf("could not create temp HOME: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpHome) })

	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(),
		"HOME="+tmpHome,
		"RAZORPAY_KEY_ID="+keyID,
		"RAZORPAY_KEY_SECRET="+keySecret,
	)
	if stdin != nil {
		cmd.Stdin = stdin
	}

	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err = cmd.Run()
	return result{stdout: out.String(), stderr: errb.String(), err: err}
}

// runJSON runs the CLI, fails the test on a non-zero exit, and parses
// the stdout as a JSON object.
func runJSON(t *testing.T, args ...string) map[string]any {
	t.Helper()
	r := run(t, args...)
	if r.err != nil {
		t.Fatalf("`razorpay %s` failed: %v\nstdout: %s\nstderr: %s",
			strings.Join(args, " "), r.err, r.stdout, r.stderr)
	}
	return parseJSON(t, args, r.stdout)
}

func parseJSON(t *testing.T, args []string, out string) map[string]any {
	t.Helper()
	out = strings.TrimSpace(out)
	var v map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("`razorpay %s` produced non-JSON output: %v\nstdout: %s",
			strings.Join(args, " "), err, out)
	}
	return v
}

// requireEntity asserts the response carries the expected Razorpay `entity` tag.
func requireEntity(t *testing.T, resp map[string]any, want string) {
	t.Helper()
	got, _ := resp["entity"].(string)
	if got != want {
		t.Fatalf("entity mismatch: want %q, got %q (resp: %v)", want, got, resp)
	}
}

// firstItem returns the first item from a `collection` response, or nil if empty.
func firstItem(t *testing.T, resp map[string]any) map[string]any {
	t.Helper()
	requireEntity(t, resp, "collection")
	items, _ := resp["items"].([]any)
	if len(items) == 0 {
		return nil
	}
	first, _ := items[0].(map[string]any)
	return first
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}
