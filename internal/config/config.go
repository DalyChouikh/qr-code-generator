// Package config provides configuration structures for QR code generation.
package config

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"
)

// OutputFormat represents the output file format.
type OutputFormat string

const (
	FormatPNG OutputFormat = "png"
	FormatSVG OutputFormat = "svg"
)

// QRConfig holds all configuration options for QR code generation.
type QRConfig struct {
	Content    string       // The URL or text to encode
	Format     OutputFormat // Output format (PNG or SVG)
	Size       int          // Dimensions in pixels (width = height)
	Foreground color.RGBA   // QR code color
	Background color.RGBA   // Background color
	OutputPath string       // Where to save the file
}

// DefaultConfig returns a QRConfig with sensible defaults.
func DefaultConfig() *QRConfig {
	return &QRConfig{
		Content:    "",
		Format:     FormatPNG,
		Size:       256,
		Foreground: color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black
		Background: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
		OutputPath: "qrcode.png",
	}
}

// Validate checks if the configuration is valid.
func (c *QRConfig) Validate() error {
	if c.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	if c.Size < 64 || c.Size > 4096 {
		return fmt.Errorf("size must be between 64 and 4096 pixels")
	}
	if c.Format != FormatPNG && c.Format != FormatSVG {
		return fmt.Errorf("format must be 'png' or 'svg'")
	}
	if c.OutputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	return nil
}

// SetOutputPath sets the output path with the correct extension.
func (c *QRConfig) SetOutputPath(path string) {
	ext := filepath.Ext(path)
	if ext == "" {
		path = path + "." + string(c.Format)
	} else {
		// Replace extension with the correct format
		path = strings.TrimSuffix(path, ext) + "." + string(c.Format)
	}
	c.OutputPath = path
}

// ParseHexColor converts a hex color string to color.RGBA.
func ParseHexColor(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color format: %s (expected 6 characters)", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid red component: %w", err)
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid green component: %w", err)
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid blue component: %w", err)
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}

// ColorToHex converts a color.RGBA to hex string.
func ColorToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// PredefinedColors provides a list of commonly used colors.
var PredefinedColors = map[string]color.RGBA{
	"Black":   {R: 0, G: 0, B: 0, A: 255},
	"White":   {R: 255, G: 255, B: 255, A: 255},
	"Red":     {R: 220, G: 53, B: 69, A: 255},
	"Green":   {R: 40, G: 167, B: 69, A: 255},
	"Blue":    {R: 0, G: 123, B: 255, A: 255},
	"Purple":  {R: 111, G: 66, B: 193, A: 255},
	"Orange":  {R: 253, G: 126, B: 20, A: 255},
	"Cyan":    {R: 23, G: 162, B: 184, A: 255},
	"Pink":    {R: 232, G: 62, B: 140, A: 255},
	"Yellow":  {R: 255, G: 193, B: 7, A: 255},
	"Teal":    {R: 32, G: 201, B: 151, A: 255},
	"Indigo":  {R: 102, G: 16, B: 242, A: 255},
}

// GetPredefinedColorNames returns the list of predefined color names.
func GetPredefinedColorNames() []string {
	return []string{
		"Black", "White", "Red", "Green", "Blue", "Purple",
		"Orange", "Cyan", "Pink", "Yellow", "Teal", "Indigo",
	}
}
