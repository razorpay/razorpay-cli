//go:build e2e
// +build e2e

// Package tests runs the Razorpay CLI as a subprocess against the live
// Razorpay API and asserts on the responses. It is gated by the `e2e`
// build tag so a plain `go test ./...` never invokes it.
//
// The suite accepts any credentials (test or live) — it does not gate on
// the key prefix. Callers are responsible for pointing it at a key whose
// account they are willing to mutate.
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
	"time"
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
			"e2e: RAZORPAY_TEST_KEY_ID / RAZORPAY_TEST_KEY_SECRET (or RAZORPAY_KEY_ID / RAZORPAY_KEY_SECRET) "+
				"must be set. Skipping.")
		os.Exit(0)
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

// runOpts is the verbose form for run. Most callers should use run /
// runWithStdin / runNoCreds instead.
type runOpts struct {
	stdin   io.Reader
	noCreds bool // when true, RAZORPAY_KEY_ID/SECRET are NOT injected
}

// run executes the CLI with credentials injected via env. Each call gets a
// fresh temp HOME so `configure` writes from one test never leak into another.
func run(t *testing.T, args ...string) result {
	t.Helper()
	return runWith(t, runOpts{}, args...)
}

func runWithStdin(t *testing.T, stdin io.Reader, args ...string) result {
	t.Helper()
	return runWith(t, runOpts{stdin: stdin}, args...)
}

// runNoCreds runs the CLI without RAZORPAY_KEY_ID / RAZORPAY_KEY_SECRET in
// the environment — used by configure subtests that assert behaviour when
// no credentials are present.
func runNoCreds(t *testing.T, stdin io.Reader, args ...string) result {
	t.Helper()
	return runWith(t, runOpts{stdin: stdin, noCreds: true}, args...)
}

func runWith(t *testing.T, opts runOpts, args ...string) result {
	t.Helper()

	tmpHome, err := os.MkdirTemp("", "razorpay-cli-home-")
	if err != nil {
		t.Fatalf("could not create temp HOME: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpHome) })

	cmd := exec.Command(binPath, args...)
	env := append(os.Environ(), "HOME="+tmpHome)
	if !opts.noCreds {
		env = append(env,
			"RAZORPAY_KEY_ID="+keyID,
			"RAZORPAY_KEY_SECRET="+keySecret,
		)
	} else {
		env = append(env, "RAZORPAY_KEY_ID=", "RAZORPAY_KEY_SECRET=")
	}
	cmd.Env = env
	if opts.stdin != nil {
		cmd.Stdin = opts.stdin
	}

	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err = cmd.Run()
	return result{stdout: out.String(), stderr: errb.String(), err: err}
}

// runJSON fails the test if the CLI exits non-zero. Used when a command
// must succeed (e.g. a `list` call, or a `create` we expect to work).
func runJSON(t *testing.T, args ...string) map[string]any {
	t.Helper()
	r := run(t, args...)
	if r.err != nil {
		t.Fatalf("`razorpay %s` failed: %v\nstdout: %s\nstderr: %s",
			strings.Join(args, " "), r.err, r.stdout, r.stderr)
	}
	return parseJSON(t, args, r.stdout)
}

// runOrSkipJSON skips (rather than fails) when the CLI exits non-zero.
// Used for API calls that may legitimately fail on a given account —
// e.g. Route account creation needs KYC fields the test account may not
// have, or downstream state may not exist yet.
func runOrSkipJSON(t *testing.T, args ...string) map[string]any {
	t.Helper()
	r := run(t, args...)
	if r.err != nil {
		t.Skipf("`razorpay %s` is not exercisable on this account: %s",
			strings.Join(args, " "), strings.TrimSpace(r.stderr))
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

// findItem scans a `collection` response for the first item whose top-level
// fields match the predicate. Returns nil if none.
func findItem(t *testing.T, resp map[string]any, pred func(map[string]any) bool) map[string]any {
	t.Helper()
	requireEntity(t, resp, "collection")
	items, _ := resp["items"].([]any)
	for _, it := range items {
		m, _ := it.(map[string]any)
		if m != nil && pred(m) {
			return m
		}
	}
	return nil
}

// strField is a small convenience to extract a top-level string field.
func strField(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

// uniqSuffix returns a high-resolution suffix suitable for keying test
// resources (receipts, references, emails) so successive runs do not collide.
func uniqSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
