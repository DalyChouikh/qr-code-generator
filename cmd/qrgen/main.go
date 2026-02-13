// QR Code Generator - A terminal-based QR code generation tool.
//
// This application provides an interactive terminal UI for generating QR codes
// with customizable colors, formats (PNG/SVG), and dimensions.
package main

import (
	"fmt"
	"os"

	"github.com/DalyChouikh/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Set via ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("qrgen %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Create the TUI program
	p := tea.NewProgram(
		ui.New(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
