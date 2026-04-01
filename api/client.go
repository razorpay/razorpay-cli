package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const baseURL = "https://api.razorpay.com/v1"

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
	if err := c.requireAuth(); err != nil {
		return nil, err
	}

	u := baseURL + path
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
