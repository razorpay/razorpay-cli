//go:build e2e
// +build e2e

package tests

import (
	"bytes"
	"strings"
	"testing"
)

// TestHelp recursively visits every command in the CLI and asserts that
// --help exits cleanly with non-empty output. The subcommand tree is
// discovered dynamically, so new commands are picked up automatically.
func TestHelp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		r := run(t, "--help")
		if r.err != nil {
			t.Fatalf("--help failed: %v\nstderr: %s", r.err, r.stderr)
		}
		if strings.TrimSpace(r.stdout) == "" {
			t.Fatalf("--help produced no stdout")
		}
	})
	walkHelp(t, nil)
}

// walkHelp recursively visits every command discovered from `--help`.
func walkHelp(t *testing.T, path []string) {
	t.Helper()
	subs := discoverSubcommands(t, path...)
	for _, sub := range subs {
		sub := sub
		next := append(append([]string{}, path...), sub)
		t.Run(strings.Join(next, "_"), func(t *testing.T) {
			args := append(append([]string{}, next...), "--help")
			r := run(t, args...)
			if r.err != nil {
				t.Fatalf("--help failed: %v\nstderr: %s", r.err, r.stderr)
			}
			if strings.TrimSpace(r.stdout) == "" {
				t.Fatalf("--help produced no stdout")
			}
			walkHelp(t, next)
		})
	}
}

// discoverSubcommands runs --help under the given path and parses the
// "Available Commands:" block. Returns nil when the command has no
// subcommands.
func discoverSubcommands(t *testing.T, parents ...string) []string {
	t.Helper()
	args := append(append([]string{}, parents...), "--help")
	r := run(t, args...)
	if r.err != nil {
		return nil
	}
	var subs []string
	inSection := false
	for _, line := range strings.Split(r.stdout, "\n") {
		trimmed := strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(trimmed, "Available Commands:") {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		// Subcommand lines start with two spaces and a name; section ends on
		// a blank line or a non-indented line.
		if trimmed == "" || !strings.HasPrefix(line, "  ") {
			break
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		// Skip cobra-injected utility commands.
		if name == "completion" || name == "help" {
			continue
		}
		subs = append(subs, name)
	}
	return subs
}

// TestValidationErrors exercises the CLI's argument-validation paths so we
// catch regressions in required flags, arg counts, and parameter parsing.
func TestValidationErrors(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		expect string
	}{
		{"payments_fetch_missing_id", []string{"payments", "fetch"}, "accepts 1 arg"},
		{"payments_capture_missing_amount", []string{"payments", "capture", "pay_dummy"}, `required flag(s) "amount"`},
		{"orders_create_missing_amount", []string{"orders", "create"}, `required flag(s) "amount"`},
		{"refunds_create_missing_id", []string{"refunds", "create"}, "accepts 1 arg"},
		// `disputes contest` still accepts `--param key=value` (most other
		// commands have migrated to typed flags), so it's the canonical
		// place to exercise ParseParams' bad-input handler.
		{"disputes_contest_bad_param", []string{"disputes", "contest", "disp_x", "--param", "noequals"}, "expected format key=value"},
		{"unknown_command", []string{"bogus"}, "unknown command"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := run(t, c.args...)
			if r.err == nil {
				t.Fatalf("expected non-zero exit, got success\nstdout: %s", r.stdout)
			}
			combined := r.stdout + r.stderr
			if !strings.Contains(combined, c.expect) {
				t.Fatalf("expected output to contain %q, got:\n%s", c.expect, combined)
			}
		})
	}
}

// TestConfigure exercises every way credentials can be supplied: both flags,
// both via stdin prompt, mixed, and the empty-input error path (which is
// only reachable when no existing credentials are present in the
// environment or on disk).
func TestConfigure(t *testing.T) {
	assertSaved := func(t *testing.T, r result) {
		t.Helper()
		if r.err != nil {
			t.Fatalf("configure failed: %v\nstdout: %s\nstderr: %s", r.err, r.stdout, r.stderr)
		}
		if !strings.Contains(r.stdout, "Credentials saved to") {
			t.Fatalf("expected success message, got:\n%s", r.stdout)
		}
	}

	t.Run("both_via_flags", func(t *testing.T) {
		r := run(t, "configure", "--key-id", "rzp_test_flagID", "--key-secret", "flag_secret")
		assertSaved(t, r)
	})

	t.Run("both_via_stdin", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("rzp_test_stdinID\nstdin_secret\n"), "configure")
		assertSaved(t, r)
	})

	t.Run("flag_id_stdin_secret", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("only_secret\n"),
			"configure", "--key-id", "rzp_test_idflag")
		assertSaved(t, r)
	})

	t.Run("flag_secret_stdin_id", func(t *testing.T) {
		r := runWithStdin(t, bytes.NewBufferString("rzp_test_idstdin\n"),
			"configure", "--key-secret", "only_secret_flag")
		assertSaved(t, r)
	})

	t.Run("empty_input_errors_with_no_existing_creds", func(t *testing.T) {
		// With no creds in the env, pressing Enter twice should error rather
		// than silently saving empty strings.
		r := runNoCreds(t, bytes.NewBufferString("\n\n"), "configure")
		if r.err == nil {
			t.Fatalf("expected error for empty input, got success\nstdout: %s", r.stdout)
		}
		if !strings.Contains(r.stdout+r.stderr, "cannot be empty") {
			t.Fatalf("expected 'cannot be empty' error, got: %s", r.stdout+r.stderr)
		}
	})
}
