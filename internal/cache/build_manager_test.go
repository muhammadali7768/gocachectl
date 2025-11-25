package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildManager_GetStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-build-stats")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mixed files
	files := map[string]string{
		"t1-d": "ok \ttest1",   // Test entry
		"b1-d": "\x7fELFbuild", // Build entry
		"b2-a": "archive",      // Build entry
		"b3":   "random",       // Build entry (no suffix, but not test)
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	mgr, err := NewBuildManager(tmpDir)
	if err != nil {
		t.Fatalf("NewBuildManager failed: %v", err)
	}

	stats, err := mgr.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	buildStats, ok := stats.(*BuildCacheStats)
	if !ok {
		t.Fatalf("Expected *BuildCacheStats, got %T", stats)
	}

	// Should find 3 build files (b1-d, b2-a, b3) and ignore t1-d
	if buildStats.EntryCount != 3 {
		t.Errorf("Expected 3 entries, got %d", buildStats.EntryCount)
	}
}

func TestBuildManager_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-build-clear")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mixed files
	files := map[string]string{
		"t1-d": "ok \ttest1",
		"b1-d": "\x7fELFbuild",
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	mgr, err := NewBuildManager(tmpDir)
	if err != nil {
		t.Fatalf("NewBuildManager failed: %v", err)
	}

	deleted, _, err := mgr.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected 1 deleted file, got %d", deleted)
	}

	// Verify b1-d is gone but t1-d remains
	if _, err := os.Stat(filepath.Join(tmpDir, "b1-d")); !os.IsNotExist(err) {
		t.Error("b1-d should have been deleted")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "t1-d")); os.IsNotExist(err) {
		t.Error("t1-d should still exist")
	}
}
