// Package generator provides QR code generation functionality.
package generator

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/DalyChouikh/internal/config"
	"github.com/skip2/go-qrcode"
)

// Generator handles QR code generation.
type Generator struct {
	config *config.QRConfig
}

// New creates a new Generator with the given configuration.
func New(cfg *config.QRConfig) *Generator {
	return &Generator{config: cfg}
}

// Generate creates the QR code and saves it to the specified path.
func (g *Generator) Generate() error {
	if err := g.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure output directory exists
	dir := filepath.Dir(g.config.OutputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	switch g.config.Format {
	case config.FormatPNG:
		return g.generatePNG()
	case config.FormatSVG:
		return g.generateSVG()
	default:
		return fmt.Errorf("unsupported format: %s", g.config.Format)
	}
}

// generatePNG creates a PNG QR code.
func (g *Generator) generatePNG() error {
	qrc, err := qrcode.New(g.config.Content, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to create QR code: %w", err)
	}

	qrc.ForegroundColor = g.config.Foreground
	qrc.BackgroundColor = g.config.Background
	qrc.DisableBorder = false

	img := qrc.Image(g.config.Size)

	file, err := os.Create(g.config.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// generateSVG creates an SVG QR code.
func (g *Generator) generateSVG() error {
	qrc, err := qrcode.New(g.config.Content, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to create QR code: %w", err)
	}

	svg := g.createSVG(qrc)

	if err := os.WriteFile(g.config.OutputPath, []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	return nil
}

// createSVG generates SVG content from a QR code.
func (g *Generator) createSVG(qrc *qrcode.QRCode) string {
	var buf bytes.Buffer

	bitmap := qrc.Bitmap()
	moduleCount := len(bitmap)

	// Calculate module size for the target dimension
	moduleSize := float64(g.config.Size) / float64(moduleCount)

	fgColor := colorToSVG(g.config.Foreground)
	bgColor := colorToSVG(g.config.Background)

	// SVG header
	buf.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" version="1.1" viewBox="0 0 %d %d" width="%d" height="%d">
`, g.config.Size, g.config.Size, g.config.Size, g.config.Size))

	// Background
	buf.WriteString(fmt.Sprintf(`  <rect width="100%%" height="100%%" fill="%s"/>
`, bgColor))

	// QR code modules
	buf.WriteString(fmt.Sprintf(`  <g fill="%s">
`, fgColor))

	for y, row := range bitmap {
		for x, module := range row {
			if module {
				px := float64(x) * moduleSize
				py := float64(y) * moduleSize
				buf.WriteString(fmt.Sprintf(`    <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f"/>
`, px, py, moduleSize, moduleSize))
			}
		}
	}

	buf.WriteString(`  </g>
</svg>`)

	return buf.String()
}

// colorToSVG converts a color.RGBA to an SVG-compatible color string.
func colorToSVG(c color.RGBA) string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.R, c.G, c.B)
}

// GetOutputPath returns the configured output path.
func (g *Generator) GetOutputPath() string {
	return g.config.OutputPath
}
