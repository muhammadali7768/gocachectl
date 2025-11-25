package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTestEntry(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gocachectl-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		filename string
		content  string
		want     bool
	}{
		{
			name:     "Valid test output (ok)",
			filename: "test1-d",
			content:  "ok \tgithub.com/example/pkg\t0.001s\n",
			want:     true,
		},
		{
			name:     "Valid test output (FAIL)",
			filename: "test2-d",
			content:  "FAIL\tgithub.com/example/pkg\t0.001s\n",
			want:     true,
		},
		{
			name:     "Valid test output (=== RUN)",
			filename: "test3-d",
			content:  "=== RUN   TestExample\n",
			want:     true,
		},
		{
			name:     "Not a test entry (no -d suffix)",
			filename: "test4",
			content:  "ok \tgithub.com/example/pkg\t0.001s\n",
			want:     false,
		},
		{
			name:     "Build artifact (ELF)",
			filename: "build1-d",
			content:  "\x7fELF\x01\x01\x01\x00",
			want:     false,
		},
		{
			name:     "Build artifact (Archive)",
			filename: "build2-d",
			content:  "!<arch>\n",
			want:     false,
		},
		{
			name:     "Build artifact (Go Object)",
			filename: "build3-d",
			content:  "go object linux/amd64",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(path, []byte(tt.content), 0644)
			if err != nil {
				t.Fatal(err)
			}

			if got := isTestEntry(path); got != tt.want {
				t.Errorf("isTestEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestManager_GetStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-stats")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some test files
	files := map[string]string{
		"t1-d": "ok \ttest1",
		"t2-d": "FAIL\ttest2",
		"b1-d": "\x7fELFbuild",
		"b2-a": "archive",
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	mgr, err := NewTestManager(tmpDir)
	if err != nil {
		t.Fatalf("NewTestManager failed: %v", err)
	}

	stats, err := mgr.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	testStats, ok := stats.(*TestCacheStats)
	if !ok {
		t.Fatalf("Expected *TestCacheStats, got %T", stats)
	}

	// Should find 2 test files (t1-d, t2-d)
	if testStats.EntryCount != 2 {
		t.Errorf("Expected 2 entries, got %d", testStats.EntryCount)
	}
}

func TestTestManager_Clear(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gocachectl-clear")
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

	mgr, err := NewTestManager(tmpDir)
	if err != nil {
		t.Fatalf("NewTestManager failed: %v", err)
	}

	deleted, _, err := mgr.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected 1 deleted file, got %d", deleted)
	}

	// Verify t1-d is gone but b1-d remains
	if _, err := os.Stat(filepath.Join(tmpDir, "t1-d")); !os.IsNotExist(err) {
		t.Error("t1-d should have been deleted")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "b1-d")); os.IsNotExist(err) {
		t.Error("b1-d should still exist")
	}
}
