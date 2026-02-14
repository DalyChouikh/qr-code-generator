// Terminal QR code preview rendering.
//
// This file provides functionality to render QR codes directly in the terminal
// using Unicode half-block characters (▀ U+2580) and ANSI color codes.
// The output is designed for maximum contrast and is scannable by QR code
// readers directly from the terminal.
//
// Compatible with modern terminals on macOS, Linux, and Windows (Windows Terminal,
// PowerShell 7+, and cmd.exe on Windows 10+).
package generator

import (
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
)

// Precomputed ANSI escape sequences for the four possible cell states.
// Using the upper half block character (▀ U+2580):
//   - Foreground color fills the TOP half of the character cell
//   - Background color fills the BOTTOM half of the character cell
//
// Index: [topIsDark][bottomIsDark]
var qrANSICodes = [2][2]string{
	// top=light
	{
		"\033[97;107m▀", // top=light, bottom=light: bright white fg + bright white bg
		"\033[97;40m▀",  // top=light, bottom=dark:  bright white fg + black bg
	},
	// top=dark
	{
		"\033[30;107m▀", // top=dark, bottom=light: black fg + bright white bg
		"\033[30;40m▀",  // top=dark, bottom=dark:  black fg + black bg
	},
}

const ansiReset = "\033[0m"

// GenerateTerminalPreview creates a terminal-renderable QR code string using
// Unicode half-block characters. Each character cell represents two vertical
// pixels, effectively doubling the vertical resolution compared to using
// full characters.
//
// The rendering uses explicit ANSI color codes (black for QR modules, bright
// white for background) to ensure consistent display regardless of terminal
// color scheme. The output includes the QR code's quiet zone (border) which
// is required for reliable scanning.
func GenerateTerminalPreview(content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	qrc, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("failed to create QR code for preview: %w", err)
	}

	bitmap := qrc.Bitmap()
	return renderBitmapToTerminal(bitmap), nil
}

// renderBitmapToTerminal converts a QR code bitmap into a terminal-renderable
// string. The bitmap is processed in pairs of rows, where each pair produces
// one line of output using the upper half block character (▀).
func renderBitmapToTerminal(bitmap [][]bool) string {
	if len(bitmap) == 0 {
		return ""
	}

	rows := len(bitmap)
	cols := len(bitmap[0])

	var buf strings.Builder
	// Estimate: ~15 bytes per cell (ANSI codes + UTF-8 char) + line overhead
	buf.Grow((rows/2 + 1) * (cols*15 + 10))

	for y := 0; y < rows; y += 2 {
		for x := 0; x < cols; x++ {
			topDark := bitmap[y][x]
			bottomDark := false
			if y+1 < rows {
				bottomDark = bitmap[y+1][x]
			}

			topIdx := boolToInt(topDark)
			bottomIdx := boolToInt(bottomDark)
			buf.WriteString(qrANSICodes[topIdx][bottomIdx])
		}
		buf.WriteString(ansiReset)
		buf.WriteString("\n")
	}

	return buf.String()
}

// boolToInt converts a boolean to an integer index (0 or 1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
