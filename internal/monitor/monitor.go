package monitor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxEvents = 500

// Entry describes one file or directory relative to the monitored root.
type Entry struct {
	Path       string    `json:"path"`
	IsDir      bool      `json:"isDir"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

// Event represents a detected filesystem change.
type Event struct {
	Type string    `json:"type"`
	Path string    `json:"path"`
	At   time.Time `json:"at"`
}

// Snapshot is the current view returned to the dashboard and API clients.
type Snapshot struct {
	Root      string    `json:"root"`
	UpdatedAt time.Time `json:"updatedAt"`
	Entries   []Entry   `json:"entries"`
	Events    []Event   `json:"events"`
	Error     string    `json:"error,omitempty"`
}

// Monitor periodically scans a directory tree and records changes.
type Monitor struct {
	root     string
	interval time.Duration

	mu        sync.RWMutex
	entries   map[string]Entry
	events    []Event
	updatedAt time.Time
	errText   string
}

func New(root string, interval time.Duration) (*Monitor, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve root: %w", err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return nil, fmt.Errorf("stat root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root %q is not a directory", absolute)
	}
	if interval < 100*time.Millisecond {
		interval = 100 * time.Millisecond
	}
	entries, err := scan(absolute)
	if err != nil {
		return nil, fmt.Errorf("scan root: %w", err)
	}
	return &Monitor{root: absolute, interval: interval, entries: entries, updatedAt: time.Now().UTC()}, nil
}

func (m *Monitor) Root() string { return m.root }
func (m *Monitor) Interval() time.Duration { return m.interval }

// Run scans until stop is closed. It is safe to call once from a goroutine.
func (m *Monitor) Run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for range ticker.C {
		m.refresh()
	}
}

func (m *Monitor) refresh() {
	current, err := scan(m.root)
	if err != nil {
		m.mu.Lock()
		m.errText = err.Error()
		m.mu.Unlock()
		return
	}
	now := time.Now().UTC()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errText = ""

	paths := make([]string, 0, len(current)+len(m.entries))
	seen := make(map[string]struct{}, len(current)+len(m.entries))
	for path := range current {
		seen[path] = struct{}{}
	}
	for path := range m.entries {
		seen[path] = struct{}{}
	}
	for path := range seen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		old, hadOld := m.entries[path]
		fresh, hasFresh := current[path]
		switch {
		case !hadOld && hasFresh:
			m.addEvent(Event{Type: "created", Path: path, At: now})
		case hadOld && !hasFresh:
			m.addEvent(Event{Type: "deleted", Path: path, At: now})
		case hadOld && hasFresh && (old.Size != fresh.Size || !old.ModifiedAt.Equal(fresh.ModifiedAt) || old.IsDir != fresh.IsDir):
			m.addEvent(Event{Type: "modified", Path: path, At: now})
		}
	}
	m.entries = current
	m.updatedAt = now
}

func (m *Monitor) addEvent(event Event) {
	m.events = append(m.events, event)
	if len(m.events) > maxEvents {
		m.events = append([]Event(nil), m.events[len(m.events)-maxEvents:]...)
	}
}

func (m *Monitor) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entries := make([]Entry, 0, len(m.entries))
	for _, entry := range m.entries {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return strings.ToLower(entries[i].Path) < strings.ToLower(entries[j].Path) })
	events := append([]Event(nil), m.events...)
	return Snapshot{Root: m.root, UpdatedAt: m.updatedAt, Entries: entries, Events: events, Error: m.errText}
}

func scan(root string) (map[string]Entry, error) {
	entries := make(map[string]Entry)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		entries[relative] = Entry{Path: relative, IsDir: d.IsDir(), Size: info.Size(), ModifiedAt: info.ModTime().UTC()}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}
