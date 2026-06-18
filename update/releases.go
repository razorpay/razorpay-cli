package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const releasesURL = "https://api.github.com/repos/razorpay/razorpay-cli/releases/latest"

const fetchTimeout = 3 * time.Second

func fetchLatestVersion(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
	if err != nil {
		return "", err
	}
	// GitHub returns 403 without a User-Agent.
	req.Header.Set("User-Agent", "razorpay-cli-update-check")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github releases endpoint returned %d", resp.StatusCode)
	}

	var body struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.TagName, nil
}
