package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/output"
)

const defaultBaseURL = "https://api.razorpay.com"

// userAgent identifies the CLI to the Razorpay API on every request.
// Format: "Razorpay-CLI/<version> (<os>/<arch>)"
var userAgent = buildUserAgent("dev")

func buildUserAgent(version string) string {
	return fmt.Sprintf("Razorpay-CLI/%s (%s/%s)", version, runtime.GOOS, runtime.GOARCH)
}

func SetUserAgentVersion(version string) {
	userAgent = buildUserAgent(version)
}

type Client struct {
	keyID     string
	keySecret string
	baseURL   string
	baseErr   error
	http      *http.Client
}

func New(keyID, keySecret string) *Client {
	base, err := resolveBaseURL()
	return &Client{
		keyID:     keyID,
		keySecret: keySecret,
		baseURL:   base,
		baseErr:   err,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
}

// resolveBaseURL returns the API base URL, honouring the RAZORPAY_BASE_URL
// override. The override is useful for pointing the CLI at a local or staging
// endpoint, but every request carries the API key in an Authorization header,
// so the override is validated and is never applied silently.
func resolveBaseURL() (string, error) {
	override := strings.TrimSpace(os.Getenv("RAZORPAY_BASE_URL"))
	if override == "" {
		return defaultBaseURL, nil
	}

	u, err := url.Parse(override)
	if err != nil {
		return "", fmt.Errorf("RAZORPAY_BASE_URL %q is not a valid URL: %w", override, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("RAZORPAY_BASE_URL %q must include a scheme and host, e.g. https://api.razorpay.com", override)
	}
	// Plain HTTP would put the API key on the wire in cleartext. Loopback is
	// allowed so the endpoint can be pointed at a local test server.
	if u.Scheme != "https" && !isLoopback(u.Hostname()) {
		return "", fmt.Errorf("RAZORPAY_BASE_URL %q uses %s; only https is accepted for non-loopback hosts, because API credentials are sent with every request", override, u.Scheme)
	}

	base := strings.TrimSuffix(override, "/")
	fmt.Fprintln(os.Stderr, "warning: RAZORPAY_BASE_URL is set — sending API requests to "+base+" instead of "+defaultBaseURL)
	return base, nil
}

func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// preflight reports whether the client is safe to send a request with: the
// endpoint must be valid and the credentials must be present.
func (c *Client) preflight() error {
	if c.baseErr != nil {
		return c.baseErr
	}
	if c.keyID == "" || c.keySecret == "" {
		return fmt.Errorf("API credentials not configured; run 'razorpay configure' or set the RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET environment variables")
	}
	return nil
}

func (c *Client) do(method, path string, body interface{}, query url.Values) ([]byte, error) {
	return c.doWithHeaders(method, path, body, query, nil)
}

func (c *Client) doWithHeaders(method, path string, body interface{}, query url.Values, extraHeaders map[string]string) ([]byte, error) {
	if err := c.preflight(); err != nil {
		return nil, err
	}

	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, reqBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.keyID, c.keySecret)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	return c.do(http.MethodGet, path, nil, query)
}

func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.do(http.MethodPost, path, body, nil)
}

func (c *Client) Patch(path string, body interface{}) ([]byte, error) {
	return c.do(http.MethodPatch, path, body, nil)
}

func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	return c.do(http.MethodPut, path, body, nil)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) GetWithHeaders(path string, query url.Values, headers map[string]string) ([]byte, error) {
	return c.doWithHeaders(http.MethodGet, path, nil, query, headers)
}

func (c *Client) PostWithHeaders(path string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.doWithHeaders(http.MethodPost, path, body, nil, headers)
}

// PostMultipart uploads a file and additional form fields using multipart/form-data.
func (c *Client) PostMultipart(path string, filePath string, fields map[string]string) ([]byte, error) {
	if err := c.preflight(); err != nil {
		return nil, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", filePath, err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, f); err != nil {
		return nil, err
	}

	for k, v := range fields {
		if err = w.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.keyID, c.keySecret)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// PrettyPrint renders the raw JSON returned by the API using the format
// the user has configured (json / yaml / toml / …). The wire format stays
// JSON; this is purely a presentation translation.
func PrettyPrint(data []byte) {
	output.Print(config.OutputFormat(), data)
}

// ParseParams parses key=value pairs from a slice of strings into a map.
func ParseParams(pairs []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid parameter %q: expected format key=value", p)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
