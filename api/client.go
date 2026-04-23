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
)

const baseURL = "https://api.razorpay.com/v1"
const apiRoot = "https://api.razorpay.com"

type Client struct {
	keyID     string
	keySecret string
	http      *http.Client
}

func New(keyID, keySecret string) *Client {
	return &Client{
		keyID:     keyID,
		keySecret: keySecret,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) requireAuth() error {
	if c.keyID == "" || c.keySecret == "" {
		return fmt.Errorf("API credentials not configured. Run 'razorpay configure' or set RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET")
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

	// Paths that start with /v2/ need the bare API root instead of /v1 base.
	base := baseURL
	if strings.HasPrefix(path, "/v2/") {
		base = apiRoot
	}
	u := base + path
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
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
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

	req, err := http.NewRequest(http.MethodPost, baseURL+path, &buf)
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

// PrettyPrint formats and prints JSON to stdout.
func PrettyPrint(data []byte) {
	var out interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		fmt.Fprintln(os.Stdout, string(data))
		return
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}

// ParseParams parses key=value pairs from a slice of strings into a map.
func ParseParams(pairs []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid parameter %q: expected key=value", p)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
