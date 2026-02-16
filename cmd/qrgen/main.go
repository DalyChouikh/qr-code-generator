// QR Code Generator - A terminal-based QR code generation tool.
//
// This application provides an interactive terminal UI for generating QR codes
// with customizable colors, formats (PNG/SVG), and dimensions.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/DalyChouikh/internal/config"
	"github.com/DalyChouikh/internal/generator"
	"github.com/DalyChouikh/internal/history"
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

		case "history":
			handleHistory()
			os.Exit(0)

		case "regen":
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: qrgen regen <id>")
				os.Exit(1)
			}
			handleRegen(os.Args[2])
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
  qrgen history         Show generation history
  qrgen regen <id>      Re-generate a QR code from history
  qrgen update          Update qrgen to the latest version
  qrgen check-update    Check if a newer version is available

Flags:
  -v, --version         Print version information
  -h, --help            Show this help message
`)
}

func handleHistory() {
	store, err := history.NewStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load history: %v\n", err)
		return
	}
	fmt.Println(store.FormatTable())
}

func handleRegen(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", idStr)
		return
	}

	store, err := history.NewStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load history: %v\n", err)
		return
	}

	entry, err := store.Get(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Reconstruct config from history entry
	fgColor, _ := config.ParseHexColor(entry.FgColor)
	bgColor, _ := config.ParseHexColor(entry.BgColor)
	cfg := &config.QRConfig{
		Content:    entry.Content,
		Format:     config.OutputFormat(entry.Format),
		Size:       entry.Size,
		Foreground: fgColor,
		Background: bgColor,
		OutputPath: entry.OutputPath,
	}

	gen := generator.New(cfg)
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Generation failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Re-generated QR code: %s\n", entry.OutputPath)
}
