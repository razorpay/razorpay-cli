// Package update prints a once-a-day notice when a newer release exists.
package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheFileName = "version-check.json"

type cache struct {
	CheckedAt time.Time `json:"checked_at"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".razorpay", cacheFileName), nil
}

// readCache returns the zero value on any error (missing, unreadable, malformed).
func readCache() cache {
	path, err := cachePath()
	if err != nil {
		return cache{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cache{}
	}
	var c cache
	if err := json.Unmarshal(data, &c); err != nil {
		return cache{}
	}
	return c
}

// writeCache writes atomically (temp file + rename) so concurrent processes
// can't tear the file.
func writeCache(c cache) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
