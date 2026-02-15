// QR Code Generator - A terminal-based QR code generation tool.
//
// This application provides an interactive terminal UI for generating QR codes
// with customizable colors, formats (PNG/SVG), and dimensions.
package main

import (
	"fmt"
	"os"

	"github.com/DalyChouikh/internal/ui"
	"github.com/DalyChouikh/internal/updater"
	tea "github.com/charmbracelet/bubbletea"
)

// Set via ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("qrgen %s (commit: %s, built: %s)\n", version, commit, date)
			os.Exit(0)

		case "update", "--update":
			fmt.Println("Checking for updates...")
			newVersion, err := updater.SelfUpdate(version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully updated to v%s!\n", newVersion)
			os.Exit(0)

		case "check-update", "--check-update":
			result, err := updater.CheckForUpdate(version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not check for updates: %v\n", err)
				os.Exit(1)
			}
			if result.UpdateAvailable {
				fmt.Printf("Update available: v%s → v%s\nRun 'qrgen update' to update.\n",
					result.CurrentVersion, result.LatestVersion)
			} else {
				fmt.Printf("You're up to date (v%s).\n", result.CurrentVersion)
			}
			os.Exit(0)

		case "--help", "-h", "help":
			printHelp()
			os.Exit(0)
		}
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

func printHelp() {
	fmt.Printf(`qrgen - A terminal-based QR code generator

Usage:
  qrgen                 Launch interactive QR code generator
  qrgen update          Update qrgen to the latest version
  qrgen check-update    Check if a newer version is available

Flags:
  -v, --version         Print version information
  -h, --help            Show this help message
`)
}
