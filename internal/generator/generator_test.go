package generator

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DalyChouikh/internal/config"
)

func newConfig(format config.OutputFormat) *config.QRConfig {
	cfg := config.DefaultConfig()
	cfg.Content = "https://example.com"
	cfg.Format = format
	return cfg
}

func TestRenderPNG(t *testing.T) {
	cfg := newConfig(config.FormatPNG)
	cfg.Size = 300

	data, err := New(cfg).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("output is not a PNG: %v", err)
	}
	if b := img.Bounds(); b.Dx() != cfg.Size || b.Dy() != cfg.Size {
		t.Errorf("PNG size = %dx%d, want %dx%d", b.Dx(), b.Dy(), cfg.Size, cfg.Size)
	}
}

func TestRenderSVG(t *testing.T) {
	cfg := newConfig(config.FormatSVG)
	cfg.Foreground = config.PredefinedColors["Red"]

	data, err := New(cfg).Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	svg := string(data)
	for _, want := range []string{"<svg", `width="256"`, `fill="rgb(220,53,69)"`} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG is missing %q", want)
		}
	}
}

func TestRenderDoesNotNeedOutputPath(t *testing.T) {
	cfg := newConfig(config.FormatPNG)
	cfg.OutputPath = ""

	if _, err := New(cfg).Render(); err != nil {
		t.Errorf("Render() error = %v, want nil", err)
	}
}

func TestRenderRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*config.QRConfig)
	}{
		{"empty content", func(c *config.QRConfig) { c.Content = "" }},
		{"size too small", func(c *config.QRConfig) { c.Size = 63 }},
		{"size too large", func(c *config.QRConfig) { c.Size = 4097 }},
		{"unknown format", func(c *config.QRConfig) { c.Format = "gif" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(config.FormatPNG)
			tt.mutate(cfg)
			if _, err := New(cfg).Render(); err == nil {
				t.Error("Render() error = nil, want an error")
			}
		})
	}
}

func TestGenerateWritesRenderedBytes(t *testing.T) {
	for _, format := range []config.OutputFormat{config.FormatPNG, config.FormatSVG} {
		t.Run(string(format), func(t *testing.T) {
			cfg := newConfig(format)
			cfg.SetOutputPath(filepath.Join(t.TempDir(), "nested", "qrcode"))

			if err := New(cfg).Generate(); err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			written, err := os.ReadFile(cfg.OutputPath)
			if err != nil {
				t.Fatalf("reading output: %v", err)
			}
			rendered, err := New(cfg).Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if !bytes.Equal(written, rendered) {
				t.Error("Generate() wrote different bytes than Render() returned")
			}
		})
	}
}

func TestGenerateRequiresOutputPath(t *testing.T) {
	cfg := newConfig(config.FormatPNG)
	cfg.OutputPath = ""

	if err := New(cfg).Generate(); err == nil {
		t.Error("Generate() error = nil, want an error")
	}
}
