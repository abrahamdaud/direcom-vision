package monitor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanAndRefreshDetectChanges(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "first.txt")
	if err := os.WriteFile(file, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := New(root, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(m.Snapshot().Entries); got != 2 {
		t.Fatalf("expected root and file entries, got %d", got)
	}
	if err := os.WriteFile(file, []byte("updated content"), 0o600); err != nil {
		t.Fatal(err)
	}
	m.refresh()
	snapshot := m.Snapshot()
	if len(snapshot.Events) != 1 || snapshot.Events[0].Type != "modified" {
		t.Fatalf("expected one modified event, got %#v", snapshot.Events)
	}
}

func TestNewRejectsFileRoot(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := New(file, time.Second); err == nil {
		t.Fatal("expected non-directory root to be rejected")
	}
}
