//go:build js && wasm

// Command qrgen-wasm exposes the qrgen generator to the website through a global `qrgen` JS object.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/DalyChouikh/internal/config"
	"github.com/DalyChouikh/internal/generator"
	"github.com/DalyChouikh/internal/templates"
)

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

func contentTypeKey(info templates.ContentTypeInfo) string {
	return strings.ToLower(info.Name)
}

func contentTypes(js.Value, []js.Value) any {
	var out []any
	for _, t := range templates.AvailableTypes() {
		out = append(out, map[string]any{
			"key":         contentTypeKey(t),
			"name":        t.Name,
			"icon":        strings.TrimSpace(t.Icon),
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
		return errorResult(fmt.Errorf("generate expects one JSON argument"))
	}

	var req generateRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return errorResult(fmt.Errorf("invalid request: %w", err))
	}

	content, err := encodeContent(req.Type, req.Fields)
	if err != nil {
		return errorResult(err)
	}

	cfg := config.DefaultConfig()
	cfg.Content = content
	cfg.Format = config.OutputFormat(req.Format)
	cfg.Size = req.Size
	if cfg.Foreground, err = config.ParseHexColor(req.FG); err != nil {
		return errorResult(fmt.Errorf("foreground: %w", err))
	}
	if cfg.Background, err = config.ParseHexColor(req.BG); err != nil {
		return errorResult(fmt.Errorf("background: %w", err))
	}

	data, err := generator.New(cfg).Render()
	if err != nil {
		return errorResult(err)
	}

	bytes := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(bytes, data)

	mime := "image/png"
	if cfg.Format == config.FormatSVG {
		mime = "image/svg+xml"
	}

	return map[string]any{"data": bytes, "mime": mime, "content": content}
}

// encodeContent mirrors the TUI form handling: fields are trimmed (except the WiFi password)
// and the same required fields are enforced before calling the template encoders.
func encodeContent(contentType string, raw map[string]string) (string, error) {
	f := make(map[string]string, len(raw))
	for k, v := range raw {
		f[k] = strings.TrimSpace(v)
	}
	f["password"] = raw["password"]

	switch contentType {
	case "url", "text":
		return f["content"], nil
	case "wifi":
		if f["ssid"] == "" {
			return "", fmt.Errorf("SSID is required")
		}
		data := templates.WiFiData{
			SSID:       f["ssid"],
			Password:   f["password"],
			Encryption: templates.WiFiEncryption(f["encryption"]),
			Hidden:     f["hidden"] == "true",
		}
		if data.Encryption == "" {
			data.Encryption = templates.WiFiWPA
		}
		return data.Encode(), nil
	case "contact":
		data := templates.VCardData{
			FirstName:    f["firstName"],
			LastName:     f["lastName"],
			Phone:        f["phone"],
			Email:        f["email"],
			Organization: f["organization"],
			Title:        f["title"],
			URL:          f["url"],
		}
		if data.FirstName == "" && data.LastName == "" {
			return "", fmt.Errorf("first or last name is required")
		}
		return data.Encode(), nil
	case "email":
		if f["address"] == "" {
			return "", fmt.Errorf("email address is required")
		}
		data := templates.EmailData{Address: f["address"], Subject: f["subject"], Body: f["body"]}
		return data.Encode(), nil
	case "sms":
		if f["phone"] == "" {
			return "", fmt.Errorf("phone number is required")
		}
		data := templates.SMSData{Phone: f["phone"], Message: f["message"]}
		return data.Encode(), nil
	default:
		return "", fmt.Errorf("unknown content type: %q", contentType)
	}
}

func errorResult(err error) any {
	return map[string]any{"error": err.Error()}
}
