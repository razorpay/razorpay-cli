package api

import (
	"strings"
	"testing"
)

func TestResolveBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		want    string
		wantErr string
	}{
		{name: "unset falls back to the default", env: "", want: defaultBaseURL},
		{name: "https override is accepted", env: "https://api-staging.razorpay.com", want: "https://api-staging.razorpay.com"},
		{name: "trailing slash is trimmed", env: "https://api-staging.razorpay.com/", want: "https://api-staging.razorpay.com"},
		{name: "loopback IPv4 may use http", env: "http://127.0.0.1:8791", want: "http://127.0.0.1:8791"},
		{name: "localhost may use http", env: "http://localhost:8791", want: "http://localhost:8791"},
		{name: "loopback IPv6 may use http", env: "http://[::1]:8791", want: "http://[::1]:8791"},

		{name: "plain http to a remote host is rejected", env: "http://example.com", wantErr: "only https"},
		{name: "non-loopback IP over http is rejected", env: "http://203.0.113.10", wantErr: "only https"},
		{name: "a bare hostname is rejected", env: "api.example.com", wantErr: "must include a scheme and host"},
		{name: "a non-http scheme is rejected", env: "ftp://example.com", wantErr: "only https"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("RAZORPAY_BASE_URL", tt.env)

			got, err := resolveBaseURL()

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected an error containing %q, got base %q and nil error", tt.wantErr, got)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected an error containing %q, got: %v", tt.wantErr, err)
				}
				if got != "" {
					t.Errorf("expected an empty base URL alongside the error, got %q", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("base URL = %q, want %q", got, tt.want)
			}
		})
	}
}

// A rejected override must surface as a request error rather than quietly
// falling back to the real API.
func TestClientRefusesToSendWithAnInvalidBaseURL(t *testing.T) {
	t.Setenv("RAZORPAY_BASE_URL", "http://example.com")

	c := New("rzp_test_key", "secret")

	if _, err := c.Get("/v1/payments", nil); err == nil {
		t.Fatal("expected Get to fail while the base URL is invalid, got nil error")
	}
	if c.baseURL == defaultBaseURL {
		t.Errorf("client silently fell back to %s instead of holding the error", defaultBaseURL)
	}
}

// Credentials are still required once the endpoint is valid.
func TestPreflightRequiresCredentials(t *testing.T) {
	t.Setenv("RAZORPAY_BASE_URL", "")

	if err := New("", "").preflight(); err == nil {
		t.Fatal("expected an error when no credentials are configured")
	}
	if err := New("rzp_test_key", "secret").preflight(); err != nil {
		t.Fatalf("unexpected error with credentials set: %v", err)
	}
}
