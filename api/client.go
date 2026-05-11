package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/output"
)

const defaultBaseURL = "https://api.razorpay.com"

type Client struct {
	keyID     string
	keySecret string
	baseURL   string
	http      *http.Client
}

func New(keyID, keySecret string) *Client {
	base := os.Getenv("RAZORPAY_BASE_URL")
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		keyID:     keyID,
		keySecret: keySecret,
		baseURL:   base,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) requireAuth() error {
	if c.keyID == "" || c.keySecret == "" {
		return fmt.Errorf("API credentials not configured; run 'razorpay configure' or set the RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET environment variables")
	}
	return nil
}

func (c *Client) do(method, path string, body interface{}, query url.Values) ([]byte, error) {
	return c.doWithHeaders(method, path, body, query, nil)
}

func (c *Client) doWithHeaders(method, path string, body interface{}, query url.Values, extraHeaders map[string]string) ([]byte, error) {
	if err := c.requireAuth(); err != nil {
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
	if err := c.requireAuth(); err != nil {
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
