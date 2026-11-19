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

	data, err := g.Render()
	if err != nil {
		return err
	}

	// Ensure output directory exists
	dir := filepath.Dir(g.config.OutputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	if err := os.WriteFile(g.config.OutputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// Render encodes the QR code in the configured format and returns the file contents.
func (g *Generator) Render() ([]byte, error) {
	if err := g.config.ValidateRendering(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	qrc, err := qrcode.New(g.config.Content, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	switch g.config.Format {
	case config.FormatPNG:
		return g.renderPNG(qrc)
	case config.FormatSVG:
		return []byte(g.createSVG(qrc)), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", g.config.Format)
	}
}

func (g *Generator) renderPNG(qrc *qrcode.QRCode) ([]byte, error) {
	qrc.ForegroundColor = g.config.Foreground
	qrc.BackgroundColor = g.config.Background
	qrc.DisableBorder = false

	var buf bytes.Buffer
	if err := png.Encode(&buf, qrc.Image(g.config.Size)); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
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
