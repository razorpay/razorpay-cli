package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// The config file holds a live API key secret, so it must not be readable by
// other users on the machine.
func TestSaveWritesConfigOwnerReadableOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file modes are not meaningful on Windows")
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := Save("rzp_test_key", "secret"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(filepath.Join(home, configDir, configFile+"."+configType))
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Errorf("config file mode = %#o, want 0600", got)
	}
}

// A config written by an older version is 0644. Saving again must narrow it,
// which WriteConfigAs alone does not do -- OpenFile only applies the mode when
// it creates the file.
func TestSaveNarrowsPermissionsOnAnExistingConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix file modes are not meaningful on Windows")
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, configFile+"."+configType)
	if err := os.WriteFile(path, []byte("key_id: old"), 0644); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	if err := Save("rzp_test_key", "secret"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Errorf("config file mode = %#o, want 0600 after re-save", got)
	}
}
