//go:build js && wasm

// Command qrgen-wasm exposes the qrgen generator to the website through a global `qrgen` JS object.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/DalyChouikh/internal/config"
	"github.com/DalyChouikh/internal/generator"
	"github.com/DalyChouikh/internal/templates"
)

// contentTypeKeys are the stable identifiers the website uses; display names may change freely.
var contentTypeKeys = map[templates.ContentType]string{
	templates.ContentURL:   "url",
	templates.ContentWiFi:  "wifi",
	templates.ContentVCard: "contact",
	templates.ContentEmail: "email",
	templates.ContentSMS:   "sms",
	templates.ContentText:  "text",
}

type generateRequest struct {
	Type   string            `json:"type"`
	Fields map[string]string `json:"fields"`
	Format string            `json:"format"`
	Size   int               `json:"size"`
	FG     string            `json:"fg"`
	BG     string            `json:"bg"`
}

func main() {
	js.Global().Set("qrgen", js.ValueOf(map[string]any{
		"generate":            js.FuncOf(generate),
		"contentTypes":        js.FuncOf(contentTypes),
		"wifiEncryptionTypes": js.FuncOf(wifiEncryptionTypes),
		"colors":              js.FuncOf(colors),
	}))

	select {}
}

func contentTypes(js.Value, []js.Value) any {
	var out []any
	for _, t := range templates.AvailableTypes() {
		out = append(out, map[string]any{
			"key":         contentTypeKeys[t.Type],
			"name":        t.Name,
			"description": t.Description,
		})
	}
	return out
}

func wifiEncryptionTypes(js.Value, []js.Value) any {
	var out []any
	for _, e := range templates.WiFiEncryptionTypes() {
		out = append(out, map[string]any{"value": string(e.Type), "name": e.Name})
	}
	return out
}

func colors(js.Value, []js.Value) any {
	var out []any
	for _, name := range config.GetPredefinedColorNames() {
		out = append(out, map[string]any{"name": name, "hex": config.ColorToHex(config.PredefinedColors[name])})
	}
	return out
}

func generate(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return errorResult(errors.New("generate expects one JSON argument"))
	}

	var req generateRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return errorResult(fmt.Errorf("invalid request: %w", err))
	}

	content, err := encodeContent(req.Type, trimFields(req.Fields))
	if err != nil {
		return errorResult(err)
	}

	cfg, err := buildConfig(req, content)
	if err != nil {
		return errorResult(err)
	}

	data, err := generator.New(cfg).Render()
	if err != nil {
		return errorResult(err)
	}

	jsData := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(jsData, data)

	return map[string]any{"data": jsData, "mime": mimeType(cfg.Format), "content": content}
}

func buildConfig(req generateRequest, content string) (*config.QRConfig, error) {
	cfg := config.DefaultConfig()
	cfg.Content = content
	cfg.Format = config.OutputFormat(req.Format)
	cfg.Size = req.Size

	var err error
	if cfg.Foreground, err = config.ParseHexColor(req.FG); err != nil {
		return nil, fmt.Errorf("foreground: %w", err)
	}
	if cfg.Background, err = config.ParseHexColor(req.BG); err != nil {
		return nil, fmt.Errorf("background: %w", err)
	}
	return cfg, nil
}

// trimFields matches the TUI, which trims every field except the WiFi password.
func trimFields(raw map[string]string) map[string]string {
	fields := make(map[string]string, len(raw))
	for name, value := range raw {
		fields[name] = strings.TrimSpace(value)
	}
	if password, ok := raw["password"]; ok {
		fields["password"] = password
	}
	return fields
}

func encodeContent(key string, fields map[string]string) (string, error) {
	if key == contentTypeKeys[templates.ContentURL] || key == contentTypeKeys[templates.ContentText] {
		return fields["content"], nil
	}

	data, err := templateFor(key, fields)
	if err != nil {
		return "", err
	}
	if err := data.Validate(); err != nil {
		return "", err
	}
	return data.Encode(), nil
}

func templateFor(key string, f map[string]string) (templates.Template, error) {
	switch key {
	case contentTypeKeys[templates.ContentWiFi]:
		return &templates.WiFiData{
			SSID:       f["ssid"],
			Password:   f["password"],
			Encryption: templates.WiFiEncryption(f["encryption"]),
			Hidden:     f["hidden"] == "true",
		}, nil
	case contentTypeKeys[templates.ContentVCard]:
		return &templates.VCardData{
			FirstName:    f["firstName"],
			LastName:     f["lastName"],
			Phone:        f["phone"],
			Email:        f["email"],
			Organization: f["organization"],
			Title:        f["title"],
			URL:          f["url"],
		}, nil
	case contentTypeKeys[templates.ContentEmail]:
		return &templates.EmailData{Address: f["address"], Subject: f["subject"], Body: f["body"]}, nil
	case contentTypeKeys[templates.ContentSMS]:
		return &templates.SMSData{Phone: f["phone"], Message: f["message"]}, nil
	default:
		return nil, fmt.Errorf("unknown content type: %q", key)
	}
}

func mimeType(format config.OutputFormat) string {
	if format == config.FormatSVG {
		return "image/svg+xml"
	}
	return "image/png"
}

func errorResult(err error) any {
	return map[string]any{"error": err.Error()}
}
