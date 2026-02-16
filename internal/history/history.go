// Package history manages QR code generation history.
//
// History entries are stored as JSON in ~/.config/qrgen/history.json.
// This allows users to review past generations and quickly re-generate
// QR codes with the same settings.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxEntries  = 50
	configDir   = "qrgen"
	historyFile = "history.json"
)

// Entry represents a single QR code generation record.
type Entry struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	Format     string    `json:"format"`
	Size       int       `json:"size"`
	FgColor    string    `json:"fg_color"`
	BgColor    string    `json:"bg_color"`
	OutputPath string    `json:"output_path"`
	CreatedAt  time.Time `json:"created_at"`
}

// Store manages the history file.
type Store struct {
	path    string
	entries []Entry
}

// NewStore creates a new history store.
func NewStore() (*Store, error) {
	path, err := historyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine history path: %w", err)
	}

	s := &Store{path: path}
	if err := s.load(); err != nil {
		// If file doesn't exist, start fresh
		s.entries = []Entry{}
	}

	return s, nil
}

// Add adds a new entry to the history.
func (s *Store) Add(entry Entry) error {
	// Assign next ID
	maxID := 0
	for _, e := range s.entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	entry.ID = maxID + 1
	entry.CreatedAt = time.Now()

	// Prepend (newest first)
	s.entries = append([]Entry{entry}, s.entries...)

	// Trim to max entries
	if len(s.entries) > maxEntries {
		s.entries = s.entries[:maxEntries]
	}

	return s.save()
}

// List returns all history entries (newest first).
func (s *Store) List() []Entry {
	return s.entries
}

// Get returns a specific entry by ID.
func (s *Store) Get(id int) (*Entry, error) {
	for _, e := range s.entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("entry #%d not found", id)
}

// Clear removes all history entries.
func (s *Store) Clear() error {
	s.entries = []Entry{}
	return s.save()
}

// FormatEntry returns a human-readable string for an entry.
func FormatEntry(e Entry) string {
	content := e.Content
	if len(content) > 50 {
		content = content[:47] + "..."
	}
	// Clean up multiline content for display
	content = strings.ReplaceAll(content, "\r\n", " ")
	content = strings.ReplaceAll(content, "\n", " ")

	return fmt.Sprintf("#%-3d  %s  %s  %dx%d  %s  %s",
		e.ID,
		e.CreatedAt.Format("2006-01-02 15:04"),
		strings.ToUpper(e.Format),
		e.Size, e.Size,
		content,
		e.OutputPath,
	)
}

// FormatTable returns a formatted table of all entries.
func (s *Store) FormatTable() string {
	if len(s.entries) == 0 {
		return "No history entries yet. Generate a QR code to get started!"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-4s  %-16s  %-4s  %-7s  %-50s  %s\n",
		"ID", "Date", "Fmt", "Size", "Content", "Output"))
	b.WriteString(strings.Repeat("─", 120) + "\n")

	for _, e := range s.entries {
		content := e.Content
		if len(content) > 50 {
			content = content[:47] + "..."
		}
		content = strings.ReplaceAll(content, "\r\n", " ")
		content = strings.ReplaceAll(content, "\n", " ")

		b.WriteString(fmt.Sprintf("#%-3d  %s  %-4s  %dx%-4d  %-50s  %s\n",
			e.ID,
			e.CreatedAt.Format("2006-01-02 15:04"),
			strings.ToUpper(e.Format),
			e.Size, e.Size,
			content,
			e.OutputPath,
		))
	}

	return b.String()
}

func historyPath() (string, error) {
	// os.UserConfigDir() returns the correct config directory per-platform:
	//   Linux:   $XDG_CONFIG_HOME or ~/.config
	//   macOS:   ~/Library/Application Support
	//   Windows: %AppData%
	configHome, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configHome, configDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(dir, historyFile), nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	return os.WriteFile(s.path, data, 0o644)
}
