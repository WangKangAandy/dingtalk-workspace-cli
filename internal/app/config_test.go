package app

import (
	"path/filepath"
	"testing"
)

func TestDefaultConfigDirUsesHomeDirectoryInOSSMode(t *testing.T) {
	homeDir := filepath.Join(t.TempDir(), "home")
	t.Setenv("HOME", homeDir)
	t.Setenv("DWS_CONFIG_DIR", "")

	got := defaultConfigDir()
	want := filepath.Join(homeDir, ".dws")
	if got != want {
		t.Fatalf("defaultConfigDir() = %q, want %q", got, want)
	}
}
