package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModManager_GetStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-mod-stats")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create module structure
	// cache/github.com/user/repo@v1.0.0/go.mod
	modulePath := filepath.Join(tmpDir, "github.com", "user", "repo@v1.0.0")
	err = os.MkdirAll(modulePath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(modulePath, "go.mod"), []byte("module github.com/user/repo"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mgr, err := NewModManager(tmpDir)
	if err != nil {
		t.Fatalf("NewModManager failed: %v", err)
	}

	stats, err := mgr.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	modStats, ok := stats.(*ModCacheStats)
	if !ok {
		t.Fatalf("Expected *ModCacheStats, got %T", stats)
	}

	if modStats.ModuleCount != 1 {
		t.Errorf("Expected 1 module, got %d", modStats.ModuleCount)
	}

	// Verify top modules contains our module
	found := false
	for _, m := range modStats.TopModules {
		if m.Path == filepath.Join("github.com", "user", "repo@v1.0.0") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected module not found in TopModules")
	}
}

func TestModManager_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-mod-clear")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file
	err = os.WriteFile(filepath.Join(tmpDir, "somefile"), []byte("data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mgr, err := NewModManager(tmpDir)
	if err != nil {
		t.Fatalf("NewModManager failed: %v", err)
	}

	deleted, _, err := mgr.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected 1 deleted file, got %d", deleted)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "somefile")); !os.IsNotExist(err) {
		t.Error("File should have been deleted")
	}
}
