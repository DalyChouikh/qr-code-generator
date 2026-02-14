// File picker component for interactive directory/file browsing.
//
// This component provides a terminal-based file browser that allows users to
// navigate directories, select existing files, or enter a new filename within
// the currently browsed directory.
//
// Compatible with macOS, Linux, and Windows filesystem conventions.
package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// fileEntry represents a single directory entry in the file picker.
type fileEntry struct {
	name  string
	isDir bool
}

// FilePicker provides an interactive file/directory browser component.
// It supports directory navigation, file selection, and inline filename input.
type FilePicker struct {
	dir        string      // Current directory being browsed
	entries    []fileEntry // Directory contents (sorted: parent first, then dirs, then files)
	cursor     int         // Currently highlighted entry index
	offset     int         // Scroll offset for long lists
	maxVisible int         // Maximum number of visible entries in the list
	err        error       // Last directory read error

	// Filename input within the currently browsed directory
	nameInput textinput.Model
	nameMode  bool // true when user is typing a new filename

	// Selection result
	selectedPath string // The full path selected by the user
	confirmed    bool   // Whether a selection has been confirmed
}

// NewFilePicker creates a new FilePicker starting at the current working directory.
// Falls back to the user's home directory if the working directory cannot be determined.
func NewFilePicker() FilePicker {
	nameInput := textinput.New()
	nameInput.Placeholder = "myqrcode"
	nameInput.CharLimit = 256
	nameInput.Width = 46

	startDir, err := os.Getwd()
	if err != nil {
		startDir, _ = os.UserHomeDir()
		if startDir == "" {
			startDir = "."
		}
	}

	fp := FilePicker{
		dir:        startDir,
		maxVisible: 12,
		nameInput:  nameInput,
	}
	fp.loadEntries()
	return fp
}

// Refresh reloads the directory contents. Call this when activating the file
// picker to ensure the listing is up-to-date.
func (fp *FilePicker) Refresh() {
	fp.loadEntries()
}

// loadEntries reads the current directory and populates the entries list.
// Directories are listed first (sorted alphabetically), followed by files
// (also sorted alphabetically). Hidden entries (starting with '.') are excluded
// to reduce clutter.
func (fp *FilePicker) loadEntries() {
	dirEntries, err := os.ReadDir(fp.dir)
	if err != nil {
		fp.err = err
		fp.entries = nil
		return
	}

	fp.err = nil
	fp.entries = make([]fileEntry, 0, len(dirEntries)+1)

	// Add parent directory navigation (except at filesystem root)
	if !isRootDir(fp.dir) {
		fp.entries = append(fp.entries, fileEntry{name: "..", isDir: true})
	}

	var dirs, files []fileEntry
	for _, e := range dirEntries {
		if strings.HasPrefix(e.Name(), ".") {
			continue // Skip hidden entries
		}
		entry := fileEntry{name: e.Name(), isDir: e.IsDir()}
		if e.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].name < dirs[j].name })
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })

	fp.entries = append(fp.entries, dirs...)
	fp.entries = append(fp.entries, files...)

	fp.cursor = 0
	fp.offset = 0
}

// isRootDir checks if the given path is a filesystem root directory.
func isRootDir(dir string) bool {
	if runtime.GOOS == "windows" {
		// Windows roots: "C:\", "D:\", etc.
		vol := filepath.VolumeName(dir)
		return vol != "" && dir == vol+string(filepath.Separator)
	}
	return dir == "/"
}

// navigateTo changes the current directory and reloads entries.
func (fp *FilePicker) navigateTo(dir string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fp.err = err
		return
	}
	fp.dir = absDir
	fp.loadEntries()
	fp.nameMode = false
}

// Update handles key messages for the file picker, routing to the
// appropriate handler based on the current mode (browsing or filename input).
func (fp *FilePicker) Update(msg tea.KeyMsg) tea.Cmd {
	if fp.nameMode {
		return fp.updateNameMode(msg)
	}
	return fp.updateBrowseMode(msg)
}

// updateBrowseMode handles key input during directory browsing.
func (fp *FilePicker) updateBrowseMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if fp.cursor > 0 {
			fp.cursor--
			if fp.cursor < fp.offset {
				fp.offset = fp.cursor
			}
		}
	case "down", "j":
		if fp.cursor < len(fp.entries)-1 {
			fp.cursor++
			if fp.cursor >= fp.offset+fp.maxVisible {
				fp.offset = fp.cursor - fp.maxVisible + 1
			}
		}
	case "enter":
		if len(fp.entries) == 0 {
			return nil
		}
		entry := fp.entries[fp.cursor]
		if entry.isDir {
			var newDir string
			if entry.name == ".." {
				newDir = filepath.Dir(fp.dir)
			} else {
				newDir = filepath.Join(fp.dir, entry.name)
			}
			fp.navigateTo(newDir)
		} else {
			// File selected — use its full path
			fp.selectedPath = filepath.Join(fp.dir, entry.name)
			fp.confirmed = true
		}
	case "n":
		// Switch to filename input mode
		fp.nameMode = true
		fp.nameInput.SetValue("")
		fp.nameInput.Focus()
		return textinput.Blink
	case "~":
		// Jump to home directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			fp.navigateTo(homeDir)
		}
	case "backspace", "delete":
		// Navigate to parent directory
		parent := filepath.Dir(fp.dir)
		if parent != fp.dir {
			fp.navigateTo(parent)
		}
	}
	return nil
}

// updateNameMode handles key input when typing a new filename.
func (fp *FilePicker) updateNameMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(fp.nameInput.Value())
		if name == "" {
			return nil // Ignore empty filenames
		}
		fp.selectedPath = filepath.Join(fp.dir, name)
		fp.confirmed = true
		return nil
	case "esc":
		fp.nameMode = false
		fp.nameInput.Blur()
		return nil
	}

	var cmd tea.Cmd
	fp.nameInput, cmd = fp.nameInput.Update(msg)
	return cmd
}

// UpdateTextInput handles non-key messages (e.g., cursor blink) for the
// filename text input. Should be called from the parent model's Update
// when this file picker is active.
func (fp *FilePicker) UpdateTextInput(msg tea.Msg) tea.Cmd {
	if fp.nameMode {
		var cmd tea.Cmd
		fp.nameInput, cmd = fp.nameInput.Update(msg)
		return cmd
	}
	return nil
}

// IsConfirmed returns true if the user has selected or entered a path.
func (fp *FilePicker) IsConfirmed() bool {
	return fp.confirmed
}

// SelectedPath returns the full path selected by the user.
func (fp *FilePicker) SelectedPath() string {
	return fp.selectedPath
}

// Reset clears the selection state so the picker can be reused.
func (fp *FilePicker) Reset() {
	fp.confirmed = false
	fp.selectedPath = ""
	fp.nameMode = false
	fp.nameInput.SetValue("")
	fp.nameInput.Blur()
}

// View renders the file picker UI using the provided styles.
func (fp *FilePicker) View(styles *Styles) string {
	var s strings.Builder

	// Current directory header — show abbreviated path with ~ for home dir
	displayDir := fp.dir
	if homeDir, err := os.UserHomeDir(); err == nil {
		if rel, err := filepath.Rel(homeDir, fp.dir); err == nil && !strings.HasPrefix(rel, "..") {
			if rel == "." {
				displayDir = "~"
			} else {
				displayDir = "~" + string(filepath.Separator) + rel
			}
		}
	}
	s.WriteString(styles.LabelFocused.Render("📂 " + displayDir))
	s.WriteString("\n\n")

	// Filename input mode
	if fp.nameMode {
		s.WriteString(styles.LabelFocused.Render("Filename:"))
		s.WriteString("\n")
		s.WriteString(styles.FocusedInput.Render(fp.nameInput.View()))
		s.WriteString("\n\n")
		s.WriteString(styles.Label.Render("Press Esc to go back to browsing"))
		return s.String()
	}

	// Error state
	if fp.err != nil {
		s.WriteString(styles.Error.Render("⚠ Cannot read directory: " + fp.err.Error()))
		s.WriteString("\n\n")
		s.WriteString(styles.Label.Render("Press Backspace to go to parent directory"))
		return s.String()
	}

	// Directory listing
	if len(fp.entries) == 0 {
		s.WriteString(styles.Label.Render("  (empty directory)"))
		s.WriteString("\n")
	} else {
		visibleEnd := fp.offset + fp.maxVisible
		if visibleEnd > len(fp.entries) {
			visibleEnd = len(fp.entries)
		}

		// Scroll indicator (top)
		if fp.offset > 0 {
			s.WriteString(styles.Label.Render("  ↑ more items above"))
			s.WriteString("\n")
		}

		for i := fp.offset; i < visibleEnd; i++ {
			entry := fp.entries[i]

			var icon string
			if entry.isDir {
				if entry.name == ".." {
					icon = "⬆ "
				} else {
					icon = "📁"
				}
			} else {
				icon = "📄"
			}

			label := fmt.Sprintf("%s %s", icon, entry.name)
			if i == fp.cursor {
				s.WriteString(styles.OptionActive.Render("▸ " + label))
			} else {
				s.WriteString(styles.Option.Render("  " + label))
			}
			s.WriteString("\n")
		}

		// Scroll indicator (bottom)
		if visibleEnd < len(fp.entries) {
			s.WriteString(styles.Label.Render("  ↓ more items below"))
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(styles.Label.Render("Press 'n' to enter a filename · '~' to go home"))

	return s.String()
}
